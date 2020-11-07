# table

<p align="center">
    <a href="https://github.com/vegarsti/table/releases"><img src="https://img.shields.io/github/release/vegarsti/table.svg" alt="Latest Release"></a>
    <a href="https://github.com/vegarsti/table/actions"><img src="https://github.com/vegarsti/table/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="http://goreportcard.com/report/github.com/vegarsti/table"><img src="http://goreportcard.com/badge/vegarsti/table" alt="Go ReportCard"></a>
</p>

```sh
> cat examples/imdb.csv
Title,Release Year,Estimated Budget
Shawshank Redemption,1994,$25 000 000
The Godfather,1972,$6 000 000
The Godfather: Part II,1974,$13 000 000
The Dark Knight,2008,$185 000 000
12 Angry Men,1957,$350 000

> go get "github.com/vegarsti/table"
> table < examples/imdb.csv
Title                   Release Year  Estimated Budget
Shawshank Redemption    1994          $25 000 000
The Godfather           1972          $6 000 000
The Godfather: Part II  1974          $13 000 000
The Dark Knight         2008          $185 000 000
12 Angry Men            1957          $350 000
```
