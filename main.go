package main

import (
	"GoMicroFeatures/models"
	"GoMicroFeatures/storage"

	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// User adalah model untuk entitas pengguna dalam database
type Users struct {
	ID       uint      `json:"id" gorm:"primaryKey:autoIncrement"`
	FullName string    `json:"fullname" gorm:"type: varchar(255)"`
	Address  string    `json:"address" gorm:"type: varchar(255)"`
	Gender   string    `json:"Gender" gorm:"type: varchar(255)"`
	UserName string    `json:"username" gorm:"type: varchar(255)"`
	Password string    `json:"password" gorm:"type: varchar(255)"`
	Articles []Article `json:"articles" gorm:"foreignKey:UserID"` // Tambahkan field untuk relasi
}

// Article adalah model untuk entitas artikel dalam database
type Article struct {
	ImgArticle  string `json:"imgarticle" gorm:"type: varchar(255)"`
	ArticleDate string `json:"articledate" gorm:"type: Date"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Description string `json:"description" gorm:"type: varchar(255)"`
	UserID      uint   `json:"user_id" gorm:"index"` // Tambahkan kolom UserID untuk menentukan relasi dengan pengguna
}

// Respository adalah struktur yang mengelola akses ke database
type Repository struct {
	DB *gorm.DB
}

// SetupRoutes mengatur rute-rute API menggunakan framework Fiber
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	// Rute untuk pengguna
	api.Post("/user", r.CreateUser)
	api.Delete("/user/:id", r.DeleteUser)
	api.Get("/user/:id", r.GetUserByID)
	api.Get("/userss", r.GetUsers)

	// Rute untuk artikel
	api.Post("/article", r.CreateArticle)
	api.Get("/article", r.GetArticles)
	api.Get("/article/:id", r.GetArticleByID)
	api.Delete("/article/:id", r.DeleteArticle)
	api.Put("/article/:id", r.UpdateArticle)

	// Rute untuk Paslons
	api.Post("/paslon", r.CreatePaslon)
	api.Get("/paslon", r.GetPaslons)
	api.Get("/paslon/:id", r.GetPaslonByID)
	api.Delete("/paslon/:id", r.DeletePaslon)
	api.Put("/paslon/:id", r.UpdatePaslon)

	// Rute untuk Partais
	api.Post("/partai", r.CreatePartai)
	api.Get("/partai", r.GetPartais)
	api.Get("/partai/:id", r.GetPartaiByID)
	api.Delete("/partai/:id", r.DeletePartai)
	api.Put("/partai/:id", r.UpdatePartai)

	// Rute untuk Vote
	api.Post("/vote", r.CreateVote)
	api.Get("/vote", r.GetVotes)
	api.Get("/vote/:id", r.GetVoteByID)
	api.Delete("/vote/:id", r.DeleteVote)
	api.Put("/vote/:id", r.UpdateVote)
}

// CreateUser membuat pengguna baru dalam database
func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := Users{}

	// Menguraikan data pengguna dari body permintaan
	err := context.BodyParser(&user)
	if err != nil {
		// Mengirim respons JSON jika parsing gagal
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	// Menambahkan pengguna baru ke dalam database
	err = r.DB.Create(&user).Error
	if err != nil {
		// Mengirim respons JSON jika pembuatan pengguna gagal
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create User"})
		return err
	}
	// Mengirim respons JSON berhasil jika tidak ada kesalahan
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "User Has been Added"})
	return nil
}

// GetUsers mendapatkan daftar semua pengguna dari database
func (r *Repository) GetUsers(context *fiber.Ctx) error {
	userModels := &[]models.Users{}

	// Mengambil semua pengguna dari database
	err := r.DB.Find(userModels).Error
	if err != nil {
		// Mengirim respons JSON jika pengambilan pengguna gagal
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Users"})
		return err
	}
	// Mengirim respons JSON berhasil dengan data pengguna
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "User fetch Successfully",
		"Data":    userModels,
	})
	return nil
}

// DeleteUser menghapus pengguna berdasarkan ID dari database
func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	userModels := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(userModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not delete user",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "User deleted Successfully",
	})
	return nil
}

// GetUserByID mendapatkan pengguna berdasarkan ID dari database
func (r *Repository) GetUserByID(context *fiber.Ctx) error {
	id := context.Params("id")
	userModels := &models.Users{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	// Menggunakan Preload untuk mengambil data artikel yang terkait dengan pengguna
	err := r.DB.Preload("Articles").Where("id = ?", id).First(userModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not get user"})
		return err
	}

	// Membuat struktur data untuk menampung informasi yang diinginkan
	var simplifiedData struct {
		ID       uint     `json:"id"`
		FullName string   `json:"fullname"`
		Address  string   `json:"address"`
		Gender   string   `json:"gender"`
		UserName string   `json:"username"`
		Articles []string `json:"articles"`
	}

	// Mengisi data ke dalam struktur data yang disederhanakan
	simplifiedData.ID = userModels.ID
	simplifiedData.FullName = userModels.FullName
	simplifiedData.Address = userModels.Address
	simplifiedData.Gender = userModels.Gender
	simplifiedData.UserName = userModels.UserName

	// Menambahkan judul artikel ke dalam struktur data yang disederhanakan
	for _, article := range userModels.Articles {
		simplifiedData.Articles = append(simplifiedData.Articles, article.Title)
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "User get Successfully",
		"data":    simplifiedData,
	})
	return nil
}

// CreateArticle membuat artikel baru dalam database
func (r *Repository) CreateArticle(context *fiber.Ctx) error {
	article := models.Articles{}

	// Menguraikan data artikel dari body permintaan
	err := context.BodyParser(&article)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	// Menambahkan artikel baru ke dalam database
	err = r.DB.Create(&article).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create Article"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Article Has been Added"})
	return nil
}

// GetArticles mendapatkan daftar semua artikel dari database
func (r *Repository) GetArticles(context *fiber.Ctx) error {
	articleModels := &[]models.Articles{}

	err := r.DB.Find(articleModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Articles"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Articles fetch Successfully",
		"Data":    articleModels,
	})
	return nil
}

// GetArticleByID mendapatkan artikel berdasarkan ID dari database
func (r *Repository) GetArticleByID(context *fiber.Ctx) error {
	id := context.Params("id")
	articleModels := &models.Articles{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	// Menggunakan Preload untuk mengambil data pengguna yang terkait
	err := r.DB.Preload("User").Where("id = ?", id).First(articleModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not get article"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Article get Successfully",
		"data":    articleModels,
	})
	return nil
}

// DeleteArticle menghapus artikel berdasarkan ID dari database
func (r *Repository) DeleteArticle(context *fiber.Ctx) error {
	articleModels := models.Articles{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(articleModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not delete article",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Article deleted Successfully",
	})
	return nil
}

// UpdateArticle mengupdate artikel berdasarkan ID dalam database
func (r *Repository) UpdateArticle(context *fiber.Ctx) error {
	id := context.Params("id")
	article := models.Articles{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not find article",
		})
		return err
	}

	// Menguraikan data artikel dari body permintaan
	err = context.BodyParser(&article)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	// Menyimpan perubahan ke dalam database
	err = r.DB.Save(&article).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not update Article"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Article updated Successfully",
		"data":    article,
	})
	return nil
}

// CreatePaslon membuat paslon baru dalam database
func (r *Repository) CreatePaslon(context *fiber.Ctx) error {
	paslon := models.Paslons{}

	err := context.BodyParser(&paslon)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Create(&paslon).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create Paslon"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Paslon Has been Added"})
	return nil
}

// GetPaslons mendapatkan daftar semua paslon dari database
func (r *Repository) GetPaslons(context *fiber.Ctx) error {
	paslonModels := &[]models.Paslons{}

	err := r.DB.Find(paslonModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Paslons"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Paslons fetch Successfully",
		"Data":    paslonModels,
	})
	return nil
}

// GetPaslonByID mendapatkan paslon berdasarkan ID dari database
func (r *Repository) GetPaslonByID(context *fiber.Ctx) error {
	id := context.Params("id")
	paslonModels := &models.Paslons{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Preload("Partais").Where("id = ?", id).First(paslonModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not get paslon"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Paslon get Successfully",
		"data":    paslonModels,
	})
	return nil
}

// DeletePaslon menghapus paslon berdasarkan ID dari database
func (r *Repository) DeletePaslon(context *fiber.Ctx) error {
	paslonModels := models.Paslons{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(paslonModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not delete paslon",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Paslon deleted Successfully",
	})
	return nil
}

// UpdatePaslon mengupdate paslon berdasarkan ID dalam database
func (r *Repository) UpdatePaslon(context *fiber.Ctx) error {
	id := context.Params("id")
	paslon := models.Paslons{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&paslon).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not find paslon",
		})
		return err
	}

	err = context.BodyParser(&paslon)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Save(&paslon).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not update Paslon"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Paslon updated Successfully",
		"data":    paslon,
	})
	return nil
}

// CreatePartai membuat partai baru dalam database
func (r *Repository) CreatePartai(context *fiber.Ctx) error {
	partai := models.Partais{}

	err := context.BodyParser(&partai)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Create(&partai).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create Partai"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Partai Has been Added"})
	return nil
}

// GetPartais mendapatkan daftar semua partai dari database
func (r *Repository) GetPartais(context *fiber.Ctx) error {
	partaiModels := &[]models.Partais{}

	err := r.DB.Find(partaiModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Partais"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Partais fetch Successfully",
		"Data":    partaiModels,
	})
	return nil
}

// GetPartaiByID mendapatkan partai berdasarkan ID dari database
func (r *Repository) GetPartaiByID(context *fiber.Ctx) error {
	id := context.Params("id")
	partaiModels := &models.Partais{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(partaiModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not get partai"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Partai get Successfully",
		"data":    partaiModels,
	})
	return nil
}

// DeletePartai menghapus partai berdasarkan ID dari database
func (r *Repository) DeletePartai(context *fiber.Ctx) error {
	partaiModels := models.Partais{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(partaiModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not delete partai",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Partai deleted Successfully",
	})
	return nil
}

// UpdatePartai mengupdate partai berdasarkan ID dalam database
func (r *Repository) UpdatePartai(context *fiber.Ctx) error {
	id := context.Params("id")
	partai := models.Partais{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&partai).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not find partai",
		})
		return err
	}

	err = context.BodyParser(&partai)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Save(&partai).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not update Partai"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Partai updated Successfully",
		"data":    partai,
	})
	return nil
}

// CreateVote membuat vote baru dalam database
func (r *Repository) CreateVote(context *fiber.Ctx) error {
	vote := models.Votes{}

	err := context.BodyParser(&vote)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Create(&vote).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create Vote"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Vote Has been Added"})
	return nil
}

// GetVotes mendapatkan daftar semua vote dari database
func (r *Repository) GetVotes(context *fiber.Ctx) error {
	voteModels := &[]models.Votes{}

	// Menggunakan Preload untuk mengambil data terkait dengan suara
	err := r.DB.Preload("User").Preload("Paslon").Find(voteModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Votes"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Votes fetch Successfully",
		"Data":    voteModels,
	})
	return nil
}

// GetVoteByID mendapatkan vote berdasarkan ID dari database
func (r *Repository) GetVoteByID(context *fiber.Ctx) error {
	id := context.Params("id")
	voteModels := &models.Votes{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(voteModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not get vote"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Vote get Successfully",
		"data":    voteModels,
	})
	return nil
}

// DeleteVote menghapus vote berdasarkan ID dari database
func (r *Repository) DeleteVote(context *fiber.Ctx) error {
	voteModels := models.Votes{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(voteModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not delete vote",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Vote deleted Successfully",
	})
	return nil
}

// UpdateVote mengupdate vote berdasarkan ID dalam database
func (r *Repository) UpdateVote(context *fiber.Ctx) error {
	id := context.Params("id")
	vote := models.Votes{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"Message": "Id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&vote).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": "could not find vote",
		})
		return err
	}

	err = context.BodyParser(&vote)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Save(&vote).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not update Vote"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"Message": "Vote updated Successfully",
		"data":    vote,
	})
	return nil
}

func main() {
	err := godotenv.Load(".env") // Memuat variabel lingkungan dari file .env
	if err != nil {
		log.Fatal(err)
	}

	// Menginisialisasi koneksi ke database
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not log the database ")
	}

	//MigrationDB
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db ")
	}
	err = models.MigrateArticles(db)
	if err != nil {
		log.Fatal("could not migrate db ")
	}
	err = models.MigratePaslons(db)
	if err != nil {
		log.Fatal("could not migrate db ")
	}
	err = models.MigratePartais(db)
	if err != nil {
		log.Fatal("could not migrate db ")
	}
	err = models.MigrateVotes(db)
	if err != nil {
		log.Fatal("could not migrate db ")
	}

	// Membuat instance Repository dengan database yang telah diinisialisasi
	r := Repository{
		DB: db,
	}
	app := fiber.New()  // Membuat instance aplikasi Fiber
	r.SetupRoutes(app)  // Menetapkan rute-rute API menggunakan Repository
	app.Listen(":8080") // Mendengarkan aplikasi di port 8080
}
