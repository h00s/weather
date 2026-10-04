// Command genplaces writes app/services/data/hr-places.tsv, the settlements
// PlacesService names locations after, from GeoNames' Croatia dump (CC BY 4.0).
//
//	cd backend && go run ./cmd/genplaces          # downloads HR.zip
//	cd backend && go run ./cmd/genplaces HR.zip   # or reads a local copy
//
// It keeps populated places (feature class P) except sections of a city (PPLX:
// they would name a point in Maksimir "Donji Bukovec" instead of Zagreb) and
// historical, abandoned or destroyed places.
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

const (
	dumpURL = "https://download.geonames.org/export/dump/HR.zip"
	output  = "app/services/data/hr-places.tsv"
)

var skipped = map[string]bool{"PPLX": true, "PPLH": true, "PPLQ": true, "PPLW": true, "PPLCH": true}

func main() {
	archive, err := readArchive(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	rows, err := places(archive)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	out.WriteString("# Croatian settlements from GeoNames (https://www.geonames.org, CC BY 4.0), written by cmd/genplaces.\n")
	out.WriteString("# name\tlatitude\tlongitude\tcounty (GeoNames admin1)\tpopulation\n")
	for _, r := range rows {
		out.WriteString(strings.Join(r, "\t") + "\n")
	}
	if err := os.WriteFile(output, out.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d places\n", output, len(rows))
}

func readArchive(args []string) ([]byte, error) {
	if len(args) > 0 {
		return os.ReadFile(args[0])
	}
	resp, err := http.Get(dumpURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", dumpURL, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// places reads HR.txt, tab-separated with the columns of
// https://download.geonames.org/export/dump/readme.txt.
func places(archive []byte) ([][]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}
	f, err := zr.Open("HR.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rows [][]string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		c := strings.Split(scanner.Text(), "\t")
		if len(c) < 15 {
			continue
		}
		name, lat, lon, class, code, admin1, population := c[1], c[4], c[5], c[6], c[7], c[10], c[14]
		if class != "P" || skipped[code] || admin1 == "" {
			continue
		}
		if population == "" {
			population = "0"
		}
		rows = append(rows, []string{name, lat, lon, admin1, population})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	slices.SortFunc(rows, func(a, b []string) int {
		return cmp.Or(strings.Compare(a[0], b[0]), strings.Compare(a[1], b[1]), strings.Compare(a[2], b[2]))
	})
	return rows, nil
}
