CREATE TABLE IF NOT EXISTS "vacancies"(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "salary" numeric(12,2)
);

CREATE TABLE IF NOT EXISTS "branches_vacancies" (
    "branch_id" INTEGER REFERENCES branches(id),
    "vacancy_id" INTEGER REFERENCES vacancies(id)
);