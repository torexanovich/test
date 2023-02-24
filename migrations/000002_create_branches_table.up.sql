CREATE TABLE IF NOT EXISTS "branches"(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "phone_numbers" VARCHAR[],
    "store_id" INTEGER REFERENCES stores(id)
);