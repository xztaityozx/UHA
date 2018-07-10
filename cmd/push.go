// Copyright © 2018 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "pushするやつ",
	Long:  `Spreadsheetsにpushするやつ`,
	Run: func(cmd *cobra.Command, args []string) {
		rj := readPushData()
		Push(rj)
	},
}

func Push(rj *PushData) {
	spreadsheetId := config.SpreadSheet.Id
	ctx := context.Background()
	client := getClient(ctx, config.SpreadSheet.CSPath)

	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatal(err)
	}

	data := []*sheets.ValueRange{
		{
			Range: fmt.Sprintf("%s%d:%s%d", rj.Column, rj.Start, rj.Column, rj.End),
			Values: [][]interface{}{
				rj.Data,
			},
			MajorDimension: "COLUMNS",
		},
	}

	reqest := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             data,
	}

	res, err := sheetService.Spreadsheets.Values.BatchUpdate(spreadsheetId, reqest).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", res)

	writeNewNextData(rj)

}

func writePushData(pd *PushData) {
	j, err := json.Marshal(pd)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(NextPath, j, 0644); err != nil {
		log.Fatal(err)
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, credentialFile string) *http.Client {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	ssconfig, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	tokFile := config.SpreadSheet.TokenPath
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(ssconfig)
		saveToken(tokFile, tok)
	}
	return ssconfig.Client(ctx, tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

type PushData struct {
	Column string        `json:"Column"`
	Data   []interface{} `json:"Data"`
	Start  int           `json:"Start"`
	End    int           `json:"End"`
}

func readPushData() *PushData {
	path := NextPath
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var writeData PushData
	if err := json.Unmarshal(b, &writeData); err != nil {
		log.Fatal(err)
	}
	return &writeData
}

func writeNewNextData(rd *PushData) {
	cur := rd.Column
	var next string
	st := rd.Start
	ed := rd.End

	if cur == "Z" {
		next = "E"
		st = ed + 2
		ed = st + 9
	} else {
		next = string([]byte(cur)[0] + 1)
	}
	wd := PushData{
		Column: next,
		Data:   []interface{}{},
		Start:  st,
		End:    ed,
	}
	path := NextPath
	jb, err := json.Marshal(wd)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(path, jb, 0644)
}

func init() {
	rootCmd.AddCommand(pushCmd)
	if _, err := os.Stat(config.SpreadSheet.CSPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.SpreadSheet.TokenPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(NextPath); err != nil {
		log.Fatal(err)
	}
}
