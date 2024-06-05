package db

import (
	"context"
	"log"

	"github.com/DuvanAlbarracin/movies_profiles/pkg/models"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Conn *pgxpool.Pool
}

func Init(url string) Handler {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalln("Error creating connection with the database:", err)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalln("Error while acquiring connection from the database pool:", err)
	}
	defer conn.Release()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln("Error pinging the database:", err)
	}

	log.Println("Database conection success!")

	createProfilesTable(pool)

	return Handler{Conn: pool}
}

func createProfilesTable(pool *pgxpool.Pool) (err error) {
	_, err = pool.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS profiles (id SERIAL PRIMARY KEY, first_name VARCHAR(50) NOT NULL, lastname VARCHAR(50), age INTEGER, gender INTEGER, city VARCHAR(50), role INTEGER NOT NULL);")
	if err != nil {
		log.Fatalln("Error creating the Users table")
		return
	}

	return nil
}

func CreateProfile(pool *pgxpool.Pool, profile *models.Profile) (err error) {
	_, err = pool.Exec(context.Background(),
		"INSERT INTO profiles(first_name, lastname, age, gender, city, role) VALUES($1, $2, $3, $4, $5, $6)",
		profile.FirstName,
		profile.LastName,
		profile.Age,
		profile.Gender,
		profile.City,
		profile.Role,
	)
	return
}

func DeleteProfile(pool *pgxpool.Pool, profile_id int64) (err error) {
	_, err = pool.Exec(context.Background(),
		"DELETE FROM profiles WHERE id = $1",
		profile_id,
	)
	return
}

func ModityProfile(pool *pgxpool.Pool, id int64, profileMap map[string]any, args []any) (err error) {
	query, err := utils.BuildUpdateQueryString("profiles", id, profileMap)
	if err != nil {
		return
	}

	_, err = pool.Exec(context.Background(), query, args...)
	return
}

func FindProfileById(pool *pgxpool.Pool, id int64) (models.Profile, error) {
	var profile models.Profile
	err := pool.QueryRow(context.Background(),
		"SELECT * FROM profiles WHERE id = $1", id).Scan(&profile.Id, &profile.FirstName,
		&profile.LastName, &profile.Age, &profile.Gender, &profile.City, &profile.Role,
	)
	return profile, err
}

func GetAllProfiles(pool *pgxpool.Pool) ([]models.Profile, error) {
	var profiles []models.Profile
	rows, _ := pool.Query(context.Background(),
		"SELECT * FROM profiles ORDER BY first_name ASC")

	for rows.Next() {
		profile := models.Profile{}
		rows.Scan(
			&profile.Id,
			&profile.FirstName,
			&profile.LastName,
			&profile.Age,
			&profile.Gender,
			&profile.City,
			&profile.Role,
		)
		profiles = append(profiles, profile)
	}

	err := rows.Err()

	return profiles, err
}
