CREATE TABLE IF NOT EXISTS "addresses"(
    "id" SERIAL PRIMARY KEY,
    "city" TEXT NOT NULL,
    "street_name" TEXT NOT NULL,
    "branch_id" INTEGER REFERENCES branches(id)
);