package main

import (
    "database/sql"    "fmt"
    "github.com/lib/pq"
    _ "github.com/lib/pq"
)
type Store struct {
    ID int    
	Name string
    Branches []*Branch
}
type Branch struct {
    ID int    
	Name string
    Phone_numbers []string    
	Store_id int
    Vacancy []*Vacancy
}
type Vacancy struct {
    ID int    
	Title string
    Branch_id int
}
type Response struct {
    Stores []*Store    
	Branches []*Branch
    Vacancies []*Vacancy
}

func check(err error) {    
	if err != nil {
        panic(err)    
	}
}

func main(){    
	connect := fmt.Sprintf("host=%s port=%d user=%s " + 
        "password=%s dbname=%s sslmode=disable",         
		"localhost", 5432, "amirkhan", "Amirkhan2005", "store1")


        
    db, err := sql.Open("postgres", connect)    
	check(err)
    defer db.Close()
    stores := []Store{
        {
            ID: 1,            
			Name: "Korzinka",
            Branches: []*Branch{                {
                    Store_id: 1,                    
					Name: "Novza",
                    Phone_numbers: []string{                        
						"+998330959595",
                        "+998333333333",                    
					},
                    Vacancy: []*Vacancy{                        
						{
                            Branch_id: 1,                            
							Title: "Farroshlik",
                        },                        
						{
                            Branch_id: 1,                            
							Title: "Buhgalter",
                        },                    
					},
                },                
				{
                    Store_id: 1,                    
					Name: "Beruniy",
                    Phone_numbers: []string{                        
						"+998330959599",
                        "+998333333388",                    
					},
                    Vacancy: []*Vacancy{
                        {                            
							Branch_id: 2,
                            Title: "Admin",                        
						},
                        {                            
							Branch_id: 2,
                            Title: "Haydovchi",                        
						},
                    },                
				},
            },        
		}, 
        {
            ID: 2,            
			Name: "Makro",
            Branches: []*Branch{                
				{
                    Store_id: 2,                    
					Name: "Sergeli",
                    Phone_numbers: []string{                        
						"+998990959599",
                        "+998903333388",                    
					},
                    Vacancy: []*Vacancy{                        
						{
                            Branch_id: 1,                            
							Title: "Qo'riqchi",
                        },                        
						{
                            Branch_id: 1,                            
							Title: "Yuk tashuvchi",
                        },                    
					},
                },                
				{
                    Store_id: 2,                    
					Name: "Qoraqamish",
                    Phone_numbers: []string{                        
						"+998880959599",
                        "+998883333388",                    
					},
                    Vacancy: []*Vacancy{                        
						{
                            Branch_id: 2,                            
							Title: "Sotuvchi",
                        },                        
						{
                            Branch_id: 2,                            
							Title: "Admin",
                        },                    
					},
                },            
			},
        },
        {            
			ID: 3,
            Name: "Havas",            
			Branches: []*Branch{
                {                    
					Store_id: 3,
                    Name: "Mirzo Ulug'bek",                    
					Phone_numbers: []string{
                        "+998990999999",                        
						"+998905555888",
                    },                    
					Vacancy: []*Vacancy{
                        {                            
							Branch_id: 1,
                            Title: "Qo'riqchi",                        
						},
                        {                            
							Branch_id: 1,
                            Title: "Kassir",                        
						},
						},                
					},
						{                    
							Store_id: 3,
							Name: "Chorsu",                    
							Phone_numbers: []string{
								"+998881263883",                        
								"+998905956081",
							},                    
							Vacancy: []*Vacancy{
								{                            
									Branch_id: 2,
									Title: "Farrosh",                        
								},
								{                            
									Branch_id: 2,
									Title: "Meneger",                        
								},
							},                
						},
					},        
				},
			}
	

	

    // CREATE TABLE 
    _, err1 := db.Exec("CREATE TABLE stores(id serial primary key, name varchar(20))") 
	check(err1)   
	 _, err2 := db.Exec("CREATE TABLE branches(id serial primary key, name varchar(20), store_id int references stores(id), phone_numbers varchar[])") 
	check(err2)
    _, err3 := db.Exec("CREATE TABLE vacancies(id serial primary key,  title varchar(30), branch_id int references branches(id))") 
	check(err3)    
	_, errbv := db.Exec("create table branches_vacancies(branch_id int references branches(id), vacancy_id int references vacancies(id))")
	check(errbv)
    

    // INSERT
    for _, x := range stores {        
		_, err := db.Exec("insert into stores(name) values($1)", x.Name)
        check(err)
        for _, i := range x.Branches{            
			_, err := db.Exec("insert into branches(name, phone_numbers, store_id) values($1, $2, $3)", i.Name, pq.Array(i.Phone_numbers), i.Store_id)
            check(err)
            for  _, j := range i.Vacancy {                
				_, err := db.Exec("insert into vacancies(title, branch_id) values($1, $2)", j.Title, j.Branch_id)
                check(err)            
			}
        }    
	}


        // GET
    myResp := Response{}
        
	storeRows, err := db.Query("SELECT id, name from stores")
    check(err)

    for storeRows.Next() {        
		store := Store{}
        err := storeRows.Scan(
            &store.ID,            
			&store.Name,
        )        
		check(err)

        branchRows, err := db.Query("SELECT id, name, phone_numbers from branches where store_id = $1", store.ID)
        check(err)  

        for branchRows.Next(){            
			branch := Branch{}

            err := branchRows.Scan(
                &branch.ID,                
				&branch.Name, 
                pq.Array(&branch.Phone_numbers),            
			)
            check(err)

            vacancyRows, err := db.Query(`SELECT v.id, v.title FROM vacancies v JOIN branches_vacancies br ON v.id = br.vacancy_id JOIN branches b ON             b.id = br.branch_id WHERE b.id = $1`, branch.ID)
            check(err)

            for vacancyRows.Next() {                
				vacancy := Vacancy{}

                err := vacancyRows.Scan(
                    &vacancy.ID,                    
					&vacancy.Title,
                )                
				check(err)
                branch.Vacancy = append(branch.Vacancy, &vacancy)
            }            
            store.Branches = append(store.Branches, &branch)    
        }

        myResp.Stores = append(myResp.Stores, &store)    
	}

	
	// PRINT
    for _, store := range myResp.Stores {        
		fmt.Println(store)
        for _, branch := range myResp.Branches {            
			fmt.Println(branch)
            for _, vacancy := range myResp.Vacancies {                
				fmt.Println(vacancy)
            }        
		}
    }
}
