package store

// var DB *gorm.DB

// func ConnectToDB() {
// 	var err error

// 	dsn := os.Getenv("DSN")
// 	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// }

type Store interface {
	User()
}

func NewStore() *Store {
	return &Store{}
}
