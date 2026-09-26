package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shafnybuilds/car_management_sys/models"
	carStore "github.com/shafnybuilds/car_management_sys/store/car"
	engineStore "github.com/shafnybuilds/car_management_sys/store/engine"
)

func TestJSONContract(t *testing.T) {
	payload := `{"name":"Civic","year":"2023","brand":"Honda","fuel_type":"Petrol","price":25000,"engine":{"engine_id":"e1f86b1a-0873-4c19-bae2-fc60329d0140","displacement":2000,"no_of_cylinders":4,"car_range":600}}`
	var request models.CarRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		t.Fatal(err)
	}
	if err := models.ValidationRequest(request); err != nil {
		t.Fatal(err)
	}
	var engineRequest models.EngineRequest
	if err := json.Unmarshal([]byte(`{"displacement":2000,"no_of_cylinders":4,"car_range":600}`), &engineRequest); err != nil {
		t.Fatal(err)
	}
	if err := models.ValidateEngineRequest(engineRequest); err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(models.Car{FuelType: request.FuelType, Engine: request.Engine})
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"id", "name", "year", "brand", "fuel_type", "engine", "price", "created_at", "updated_at"} {
		if _, ok := response[key]; !ok {
			t.Errorf("response is missing %q: %s", key, body)
		}
	}
	var engine map[string]json.RawMessage
	if err := json.Unmarshal(response["engine"], &engine); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"engine_id", "displacement", "no_of_cylinders", "car_range"} {
		if _, ok := engine[key]; !ok {
			t.Errorf("engine response is missing %q", key)
		}
	}
}

// A small database/sql driver exercises scanning and statement errors without a server.
// Schema initialization has an opt-in PostgreSQL check below.
type statementDB struct {
	values  []driver.Value
	err     error
	queries []string
}

func (s *statementDB) Connect(context.Context) (driver.Conn, error) { return statementConn{s}, nil }
func (s *statementDB) Driver() driver.Driver                        { return statementDriver{} }

type statementDriver struct{}

func (statementDriver) Open(string) (driver.Conn, error) { return nil, errors.New("use connector") }

type statementConn struct{ state *statementDB }

func (statementConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (statementConn) Close() error { return nil }
func (statementConn) Begin() (driver.Tx, error) {
	return nil, errors.New("single statements must not defer commit")
}
func (c statementConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	c.state.queries = append(c.state.queries, query)
	if strings.HasPrefix(query, "SELECT id FROM engine") {
		return &statementRows{values: []driver.Value{engineID}}, nil
	}
	if c.state.err != nil {
		return nil, c.state.err
	}
	return &statementRows{values: c.state.values}, nil
}
func (c statementConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	if c.state.err != nil {
		return nil, c.state.err
	}
	return driver.RowsAffected(1), nil
}

type statementRows struct {
	values []driver.Value
	done   bool
}

func (r *statementRows) Columns() []string { return make([]string, len(r.values)) }
func (*statementRows) Close() error        { return nil }
func (r *statementRows) Next(dest []driver.Value) error {
	if r.done || r.values == nil {
		return io.EOF
	}
	copy(dest, r.values)
	r.done = true
	return nil
}

const carID = "c7c1a6d5-1ec4-4c64-a59a-8a2f6f3d2bf3"
const engineID = "e1f86b1a-0873-4c19-bae2-fc60329d0140"

func TestCarReadsAndDeletion(t *testing.T) {
	base := []driver.Value{carID, "Civic", "2023", "Honda", "Petrol", engineID, "25000.00", time.Now(), time.Now()}
	for _, mode := range []string{"id", "brand", "brand-engine", "delete"} {
		t.Run(mode, func(t *testing.T) {
			values := append([]driver.Value(nil), base...)
			if mode == "id" || mode == "brand-engine" {
				values = append(values, engineID, int64(2000), int64(4), int64(600))
			}
			state := &statementDB{values: values}
			db := sql.OpenDB(state)
			defer db.Close()
			store := carStore.New(db)
			var car models.Car
			var err error
			switch mode {
			case "id":
				car, err = store.GetCarById(context.Background(), carID)
			case "delete":
				car, err = store.DeleteCar(context.Background(), carID)
			default:
				var cars []models.Car
				cars, err = store.GetCarByBrand(context.Background(), "Honda", mode == "brand-engine")
				if err == nil && len(cars) != 1 {
					t.Fatalf("got %d cars", len(cars))
				}
				if err == nil {
					car = cars[0]
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			if car.FuelType != "Petrol" || car.Engine.EngineID.String() != engineID {
				t.Fatalf("incorrect scan: %+v", car)
			}
			if (mode == "id" || mode == "brand-engine") && car.Engine.NoOfCylinders != 4 {
				t.Fatalf("engine details were discarded: %+v", car.Engine)
			}
			query := state.queries[len(state.queries)-1]
			if strings.Contains(query, "furl_type") || strings.Contains(query, "engie_id") {
				t.Fatalf("invalid column in %s", query)
			}
			if mode == "brand-engine" && !strings.Contains(query, "LEFT JOIN engine e ON") {
				t.Fatalf("missing engine alias: %s", query)
			}
			if mode == "delete" && !strings.HasPrefix(query, "DELETE FROM car") {
				t.Fatal("deletion must return the row from the DELETE statement")
			}
		})
	}
}

func TestStoreDatabaseErrors(t *testing.T) {
	failure := errors.New("database statement failed")
	state := &statementDB{err: failure}
	db := sql.OpenDB(state)
	defer db.Close()
	car := carStore.New(db)
	engine := engineStore.New(db)
	ctx := context.Background()
	checks := map[string]func() error{
		"create-car":    func() error { _, err := car.CreateCar(ctx, &models.CarRequest{}); return err },
		"update-car":    func() error { _, err := car.UpdateCar(ctx, carID, &models.CarRequest{}); return err },
		"delete-car":    func() error { _, err := car.DeleteCar(ctx, carID); return err },
		"create-engine": func() error { _, err := engine.EngineCreate(ctx, &models.EngineRequest{}); return err },
		"update-engine": func() error { _, err := engine.EngineUpdate(ctx, engineID, &models.EngineRequest{}); return err },
		"delete-engine": func() error { _, err := engine.DeleteEngine(ctx, engineID); return err },
	}
	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			if err := check(); !errors.Is(err, failure) {
				t.Fatalf("got %v, want original database error", err)
			}
		})
	}
}

func TestDeletionOfMissingRows(t *testing.T) {
	db := sql.OpenDB(&statementDB{})
	defer db.Close()
	if _, err := carStore.New(db).DeleteCar(context.Background(), carID); err == nil {
		t.Fatal("missing car deletion must report not found")
	}
	engine, err := engineStore.New(db).DeleteEngine(context.Background(), engineID)
	if err != nil || engine.EngineID != uuid.Nil {
		t.Fatalf("missing engine deletion must preserve not-found behavior: %+v, %v", engine, err)
	}
}
func TestSchemaPreservesDataPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run the PostgreSQL schema check")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	schemaName := "regression_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := tx.Exec("CREATE SCHEMA " + schemaName + "; SET LOCAL search_path TO " + schemaName); err != nil {
		t.Fatal(err)
	}
	schema, err := os.ReadFile("store/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec("INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES ($1, 2000, 4, 600)", engineID); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec("INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price) VALUES ($1, 'Civic', '2023', 'Honda', 'Petrol', $2, 25000)", carID, engineID); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM car WHERE id = $1", carID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("schema initialization reset existing car data")
	}
	if _, err := tx.Exec("DELETE FROM engine WHERE id = $1", engineID); err != nil {
		t.Fatal(err)
	}
	if err := tx.QueryRow("SELECT COUNT(*) FROM car WHERE id = $1", carID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("engine foreign key cascade is missing")
	}
}
