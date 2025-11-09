#!/usr/bin/env bash
set -e

POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
PGPASSWORD=$POSTGRES_PASSWORD
export PGPASSWORD

echo "Create cars table"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "cars" <<-EOSQL
    CREATE TABLE IF NOT EXISTS cars
    (
        id                  SERIAL PRIMARY KEY,
        car_uid             uuid UNIQUE NOT NULL,
        brand               VARCHAR(80) NOT NULL,
        model               VARCHAR(80) NOT NULL,
        registration_number VARCHAR(20) NOT NULL,
        power               INT,
        price               INT         NOT NULL,
        type                VARCHAR(20) CHECK (type IN ('SEDAN', 'SUV', 'MINIVAN', 'ROADSTER')),
        availability        BOOLEAN     NOT NULL
    );

    INSERT INTO cars (car_uid, brand, model, registration_number, power, price, type, availability) VALUES
    (gen_random_uuid(), 'Toyota', 'Camry', 'A123BC777', 150, 25000, 'SEDAN', true),
    (gen_random_uuid(), 'BMW', 'X5', 'B456DE777', 230, 45000, 'SUV', true),
    (gen_random_uuid(), 'Mercedes', 'Vito', 'C789FG777', 140, 35000, 'MINIVAN', false),
    (gen_random_uuid(), 'Ferrari', '488 GTB', 'D012HI777', 670, 250000, 'ROADSTER', true),
    (gen_random_uuid(), 'Honda', 'Civic', 'E345JK777', 120, 20000, 'SEDAN', true),
    (gen_random_uuid(), 'Jeep', 'Wrangler', 'F678LM777', 270, 40000, 'SUV', true),
    (gen_random_uuid(), 'Ford', 'Transit', 'G901NP777', 130, 30000, 'MINIVAN', false),
    (gen_random_uuid(), 'Porsche', '911', 'H234QR777', 450, 150000, 'ROADSTER', true),
    (gen_random_uuid(), 'Nissan', 'Altima', 'I567ST777', 180, 22000, 'SEDAN', true),
    (gen_random_uuid(), 'Land Rover', 'Discovery', 'J890UV777', 340, 60000, 'SUV', true)
    ON CONFLICT DO NOTHING;
EOSQL

echo "Create rentals table"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "rentals" <<-EOSQL
    CREATE TABLE IF NOT EXISTS rental
    (
        id          SERIAL PRIMARY KEY,
        rental_uid  uuid UNIQUE              NOT NULL,
        username    VARCHAR(80)              NOT NULL,
        payment_uid uuid                     NOT NULL,
        car_uid     uuid                     NOT NULL,
        date_from   TIMESTAMP WITH TIME ZONE NOT NULL,
        date_to     TIMESTAMP WITH TIME ZONE NOT NULL,
        status      VARCHAR(20)              NOT NULL CHECK (status IN ('IN_PROGRESS', 'FINISHED', 'CANCELED'))
    );
EOSQL

echo "Create payments table"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "payments" <<-EOSQL
    CREATE TABLE IF NOT EXISTS payment
    (
        id          SERIAL PRIMARY KEY,
        payment_uid uuid        NOT NULL,
        status      VARCHAR(20) NOT NULL CHECK (status IN ('PAID', 'CANCELED')),
        price       INT         NOT NULL
    );
EOSQL