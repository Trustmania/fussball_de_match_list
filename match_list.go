package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "github.com/PuerkitoBio/goquery"
)

type Match struct {
    Start_at string `json:"start_at"`
    Team_one string `json:"team_one"`
    Team_two string `json:"team_two"`
}

type MatchList struct {
    Matches []Match `json:"match"`
}

func MatchListHandler(w http.ResponseWriter, r *http.Request) {
  url := r.URL.Query().Get("url")

  doc, err := goquery.NewDocument(url)
  if err != nil {
    //log.Fatal(err)
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, err)
    return
  }

  var headline string
  var match_list MatchList

  lines := doc.Find("#team-matchplan-table tbody tr")
  lines.Each(func(i int, s *goquery.Selection) {
    if i == 0 { // headline
        headline = s.Find("td").Text()
    }

    if i == 2 { // team names
        m := Match{
            Start_at: headline,
            Team_one: s.Find("td.column-club .club-name").Text(),
            Team_two: s.Find("td.column-club .club-name").Text(),
        }
        match_list.Match = append(match_list.Match, m)
    }
  })
  fmt.Println(match_list)

  resp, _ := json.MarshalIndent(match_list, "", "  ")

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  fmt.Fprint(w, string(resp))
}