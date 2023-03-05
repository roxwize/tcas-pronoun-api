/*
	tcas-pronoun-api, RESTful API for fetching and modifying pronouns of TCaS users
	Copyright (C) 2023  thekifake

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Error = color.New(color.FgRed)
	Success = color.New(color.FgHiGreen)
	LogoBlue = color.New(color.FgBlue)
	LogoPurple = color.New(color.FgRed)
)

func welcomeText() {
	logo, err := os.ReadFile("logo.txt")
	if err == nil {
		fmt.Printf("%s\n\n",string(logo))
	}
	// Using Printf with Color.Sprint is not supported on default Windows terminal
	fmt.Print("\t\t\t=[ || This is the ")
	LogoBlue.Print("Two Cans")
	fmt.Print(" and ")
	LogoPurple.Print("String")
	fmt.Print(" pronoun API || ]=\n\n")
	fmt.Println()
}

type User struct{
	Username	string	`json:"username"`
	Pronouns	string	`json:"pronouns"`
	LineNum		int			`json:"-"`
}
var Pronouns = make([]User, 0)

func fatal(m string, e error) {
	Error.Printf("There was an error %s:\n\t%s\n", m, e)
	os.Exit(1)
}
func serverFatal(m string, c *gin.Context, e error) {
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": m})
	}
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func readPronounsFile() ([]byte, error) {
	f, err := os.ReadFile("pronouns")
	return f, err
}
func parsePronounsFile(f []byte) {
	Pronouns = make([]User, 0)
	ls := strings.Split(string(f), "\n")
	for i, l := range ls {
		s := strings.Split(l, ";")
		Pronouns = append(Pronouns, User{Username: s[0], Pronouns: s[1], LineNum: i})
	}
}

func getAllPronouns(c *gin.Context) {
	c.JSON(http.StatusOK, Pronouns)
}
func getPronoun(c *gin.Context) {
	id := c.Param("username")
	for _, u := range Pronouns {
		if u.Username == id {
			c.Data(http.StatusOK, "text/plain", []byte(u.Pronouns))
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
func addPronoun(c *gin.Context) {
	// This should require authentication
	formUsername := c.Request.FormValue("username")
	formPronouns := c.Request.FormValue("pronouns")
	for _, u := range Pronouns {
		if u.Username == formUsername {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
	}
	entry := User{Username: formUsername, Pronouns: formPronouns}
	f, err := os.ReadFile("pronouns")
	serverFatal("reading pronouns failed", c, err)
	newList := []byte(fmt.Sprintf("%s\n%s;%s", string(f), formUsername, formPronouns))
	err = os.WriteFile("pronouns", newList, os.ModeAppend)
	serverFatal("adding entry failed", c, err)
	parsePronounsFile(newList)
	c.JSON(http.StatusCreated, entry)
}
func setPronoun(c *gin.Context) {
	// This should require authentication
	formUsername := c.Param("username")
	formPronouns := c.Request.FormValue("pronouns")
	for _, u := range Pronouns {
		if u.Username == formUsername {
			if u.Pronouns == formPronouns {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}
			f, err := os.ReadFile("pronouns")
			serverFatal("reading pronouns failed", c, err)
			lines := strings.Split(string(f), "\n")
			lines[u.LineNum] = fmt.Sprintf("%s;%s", formUsername, formPronouns)
			os.WriteFile("pronouns", []byte(strings.Join(lines, "\n")), os.ModeAppend)
			f, err = os.ReadFile("pronouns")
			serverFatal("reading pronouns failed", c, err)
			parsePronounsFile(f)
			c.JSON(http.StatusOK, Pronouns[u.LineNum])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
func deletePronoun(c *gin.Context) {
	// This should require authentication
	username := c.Param("username")
	for _, u := range Pronouns {
		if u.Username == username {
			f, err := os.ReadFile("pronouns")
			serverFatal("reading pronouns failed", c, err)
			lines := strings.Split(string(f), "\n")
			lines = remove(lines, u.LineNum)
			os.WriteFile("pronouns", []byte(strings.Join(lines, "\n")), os.ModeAppend)
			f, err = os.ReadFile("pronouns")
			serverFatal("reading pronouns failed", c, err)
			parsePronounsFile(f)
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func main() {
	welcomeText()

	fmt.Println("Attempting to read from the pronouns database...")
	// TODO replace this with a call to the actual database
	f, err := readPronounsFile()
	if err != nil {
		fatal("accessing the database", err)
	}
	parsePronounsFile(f)
	Success.Println("Pronouns loaded!")
	
	router := gin.Default()
	router.GET		("/pronouns", getAllPronouns)
	router.GET		("/pronouns/:username", getPronoun)
	router.POST		("/pronouns/add", addPronoun)
	router.PATCH	("/pronouns/:username", setPronoun)
	router.DELETE	("/pronouns/:username", deletePronoun)

	router.Run(":1337")
}

