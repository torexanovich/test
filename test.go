package main

import (
	"database/sql"
	"fmt"

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
    PhoneNumbers []string    
    Address *Address
    Vacancies []*Vacancy
}

type Vacancy struct {
    ID int    
	Name string
    Salary float64
}

type Address struct {
    ID int
    City string
    StreetName string
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
		"localhost", 5432, "amirkhan", "Amirkhan2005", "migration")
        
    db, err := sql.Open("postgres", connect)    
	check(err)
    defer db.Close()

    tx, err := db.Begin()
    check(err)
   
    
    stores := []Store{
        {           
			Name: "Korzinka",
            Branches: []*Branch{                
                {                   
					Name: "Kohinur",
                    PhoneNumbers: []string{                        
						"+998330959595",
                        "+998333333333",                    
					},
                    Address: &Address{
                        City: "Toshkent Olmazor",
                        StreetName: "Beruniy",
                    },
                    Vacancies: []*Vacancy{                        
						{                            
							Name: "Kassir",
                            Salary: 4000000,
                        },                        
						{
							Name: "Oxranik",
                            Salary: 3500000,

                        },                    
					},
                },                
				{                    
					Name: "Qatorto'l",
                    PhoneNumbers: []string{                        
						"+998956545665",
                        "+998915545445",                    
					},
                    Address: &Address{
                        City: "Toshkent Chilonzor",
                        StreetName: "Mirzo Ulug'bek",
                    },
                    Vacancies: []*Vacancy{
                        {                            
                            Name: "Tozalovchi", 
                            Salary: 3000000,                       
						},
                        {                            
                            Name: "Manager",
                            Salary: 4500000,
                        
						},
                    },                
				},
            },        
		}, 
        {
			Name: "Makro",
            Branches: []*Branch{                
				{
					Name: "Sergeli",
                    PhoneNumbers: []string{                        
						"+998990959599",
                        "+998903333388",                    
					},
                    Address: &Address{
                        City: "Toshkent Sergeli",
                        StreetName: "Ozodlik ko'chasi",
                    },
                    Vacancies: []*Vacancy{                        
						{
							Name: "Kassir",
                            Salary: 3500000,

                        },                        
						{
							Name: "Oxranik",
                            Salary: 2500000,

                        },                    
					},
                },                
				{
					Name: "Makro Tinchlik metro",
                    PhoneNumbers: []string{                        
						"+998880959599",
                        "+998883333388",                    
					},
                    Address: &Address{
                        City: "Toshkent Olmazor",
                        StreetName: "Olmazor ko'chasi",
                    },
                    Vacancies: []*Vacancy{                        
						{
							Name: "Tozalovchi",
                            Salary: 3500000,

                        },                        
						{
							Name: "Admin",
                            Salary: 4000000,

                        },                    
					},
                },            
			},
        },
    } 

    

    // INSERT
    for _, store := range stores {      
        var storeID int  
		err := tx.QueryRow("INSERT INTO stores (name) VALUES ($1) RETURNING id", store.Name).Scan(&storeID)
        if err != nil {
            tx.Rollback()
            return
        }
        for _, branch := range store.Branches{  
            var branchID int  
          
            err := tx.QueryRow("INSERT INTO branches (name, phone_numbers, store_id) VALUES ($1, $2, $3) RETURNING id", branch.Name, pq.Array(branch.PhoneNumbers), storeID).Scan(&branchID)
            if err != nil {
                tx.Rollback()
                return
            }

            _, err = tx.Exec("INSERT INTO addresses(city, stree_name, branch_id) VALUES($1, $2, $3)", branch.Address.City, branch.Address.StreetName, branchID)

            for  _, vacancy := range branch.Vacancies {           
                var vacancyID int     

                err := tx.QueryRow("INSERT INTO vacancies(name, salary) VALUES ($1, $2) RETURNING id", vacancy.Name, vacancy.Salary).Scan(&vacancyID)
                if err != nil {
                    tx.Rollback()
                    return
                }      
                
                _, err = tx.Exec("INSERT INTO branches_vacancies(branch_id, vacancy_id) VALUES($1, $2)", branchID, vacancyID)
                if err != nil {
                    tx.Rollback()
                    return
                }      
                
			}
        }    
	}


    // // GET
    // myResp := Response{}
        
	// storeRows, err := db.Query("SELECT id, name from stores")
    // check(err)

    // for storeRows.Next() {        
	// 	store := Store{}
    //     err := storeRows.Scan(
    //         &store.ID,            
	// 		&store.Name,
    //     )        
	// 	check(err)

    //     branchRows, err := db.Query("SELECT id, name, phone_numbers from branches where store_id = $1", store.ID)
    //     check(err)  

    //     for branchRows.Next(){            
	// 		branch := Branch{}

    //         err := branchRows.Scan(
    //             &branch.ID,                
	// 			&branch.Name, 
    //             pq.Array(&branch.Phone_numbers),            
	// 		)
    //         check(err)

    //         vacancyRows, err := db.Query(`SELECT v.id, v.title FROM vacancies v JOIN branches_vacancies br ON v.id = br.vacancy_id JOIN branches b ON             b.id = br.branch_id WHERE b.id = $1`, branch.ID)
    //         check(err)

    //         for vacancyRows.Next() {                
	// 			vacancy := Vacancy{}

    //             err := vacancyRows.Scan(
    //                 &vacancy.ID,                    
	// 				&vacancy.Title,
    //             )                
	// 			check(err)
    //             branch.Vacancy = append(branch.Vacancy, &vacancy)
    //         }            
    //         store.Branches = append(store.Branches, &branch)    
    //     }

    //     myResp.Stores = append(myResp.Stores, &store)    
	// }

	
	// PRINT
    // for _, store := range myResp.Stores {        
	// 	fmt.Println(store)
    //     for _, branch := range myResp.Branches {            
	// 		fmt.Println(branch)
    //         for _, vacancy := range myResp.Vacancies {                
	// 			fmt.Println(vacancy)
    //         }        
	// 	}
    // }
}
