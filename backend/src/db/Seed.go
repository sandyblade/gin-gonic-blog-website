/**
 * This file is part of the Sandy Andryanto Blog Applicatione.
 *
 * @author     Sandy Andryanto <sandy.andryanto.blade@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

package db

import (
	helpers "api/backend/src/helpers"
	models "api/backend/src/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"github.com/Pallinder/go-randomdata"
	"github.com/bxcodec/faker/v4"
	"github.com/kristijorgji/goseeder"
	math "math/rand"
)

func init() {
	goseeder.Register(CreateUser)
}

func CreateUser(s goseeder.Seeder) {

	var totalRow int64

	_db := SetupDB()
	_db.Model(&models.User{}).Where("id <> 0").Count(&totalRow)

	if totalRow == 0 {
		for i := 1; i <= 10; i++ {

			bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
			if _, err := rand.Read(bytes); err != nil {
				panic(err.Error())
			}

			key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
			encrypted := helpers.Encrypt("P@ssw0rd!123", key)

			min := 1
			max := 2
			gender := math.Intn(max-min+1) + min
			firstName := ""
			genderChar := ""

			if gender == 1 {
				genderChar = "M"
				firstName = faker.FirstNameMale()
			} else {
				genderChar = "F"
				firstName = faker.FirstNameFemale()
			}

			user := models.User{
				Email:     randomdata.Email(),
				Phone:     randomdata.PhoneNumber(),
				FirstName: sql.NullString{String: firstName, Valid: true},
				LastName:  sql.NullString{String: faker.LastName(), Valid: true},
				Gender:    sql.NullString{String: genderChar, Valid: true},
				Country:   sql.NullString{String: randomdata.Country(randomdata.FullCountry), Valid: true},
				Address:   sql.NullString{String: randomdata.Address(), Valid: true},
				AboutMe:   sql.NullString{String: randomdata.Paragraph(), Valid: true},
				JobTitle:  sql.NullString{String: randomdata.SillyName(), Valid: true},
				Twitter:   sql.NullString{String: faker.Username(), Valid: true},
				Facebook:  sql.NullString{String: faker.Username(), Valid: true},
				Instagram: sql.NullString{String: faker.Username(), Valid: true},
				LinkedIn:  sql.NullString{String: faker.Username(), Valid: true},
				Salt:      key,
				Password:  encrypted,
				Confirmed: 1,
			}
			_db.Create(&user)
		}
	}

}
