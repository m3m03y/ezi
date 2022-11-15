package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/exp/slices"
)

func save_file(name, data string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		fmt.Printf("Error: %s", err2.Error())
	}
}

func save_json_file(name string, data any) {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		_ = ioutil.WriteFile(name, json_data, 0644)
	}
}

func sort_map_by_value(input_map map[string]int) map[string]int {
	keys := make([]string, 0, len(input_map))
	for key := range input_map {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return input_map[keys[i]] > input_map[keys[j]]
	})

	output_map := make(map[string]int)

	for _, key := range keys {
		output_map[key] = input_map[key]
		fmt.Printf("%-7v %v\n", key, input_map[key])
	}
	return output_map
}

func sort_map_by_float_value(input_map map[string]float64) map[string]float64 {
	keys := make([]string, 0, len(input_map))
	for key := range input_map {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return input_map[keys[i]] > input_map[keys[j]]
	})

	output_map := make(map[string]float64)

	for _, key := range keys {
		output_map[key] = input_map[key]
		fmt.Printf("%-7v %v\n", key, input_map[key])
	}
	return output_map
}

func collect_data(site_url string, allowed_domains string) map[string][]string {
	site_map := make(map[string][]string)
	visited_site_count := 0
	site_root := ""

	c := colly.NewCollector(
		// Visit only domains allowed domains
		colly.AllowedDomains(allowed_domains),
		colly.MaxDepth(3),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		temp_list := site_map[site_root]
		current_link := e.Request.AbsoluteURL(link)
		if !slices.Contains(temp_list, current_link) && (current_link != site_root) {
			if strings.Contains(current_link, allowed_domains) {
				site_map[site_root] = append(temp_list, current_link)
			}
		}
		if visited_site_count < 30 {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		visited_site_count++
		site_root = r.URL.String()
		site_map[site_root] = []string{}
	})

	c.Visit(site_url)
	fmt.Println("Visited ", visited_site_count, " sites")
	return site_map
}

func main() {
	site_map := collect_data("https://flyingwildhog.com/", "flyingwildhog.com")
	// task 1: print site map
	save_json_file("site_list.json", site_map)

	// task 1: create site list
	site_list := []string{}
	site_indexes := make(map[string]int)
	idx := 0
	for key := range site_map {
		site_list = append(site_list, key)
		site_indexes[key] = idx
		idx++
	}

	sort.Strings(site_list)

	out_string := ""
	for _, data := range site_list {
		out_string += data + "\n"
	}

	save_file("site_list.txt", out_string)

	// task 2: create adjacency matrix for site_map
	adj_matrix_txt := ""
	adj_matrix := [][]int{}
	for i := 0; i < len(site_list); i++ {
		adj_matrix_txt += fmt.Sprintf(",%d", i)
	}
	adj_matrix_txt += "\n"
	for i := 0; i < len(site_list); i++ {
		column := site_list[i]
		adj_matrix_txt += fmt.Sprintf("%d", i)
		column_values := []int{}
		for j := 0; j < len(site_list); j++ {
			row := site_list[j]
			if slices.Contains(site_map[column], row) {
				adj_matrix_txt += ",1"
				column_values = append(column_values, 1)
			} else {
				adj_matrix_txt += ",0"
				column_values = append(column_values, 0)
			}
		}
		adj_matrix_txt += "\n"
		adj_matrix = append(adj_matrix, column_values)
	}

	save_file("adj_matrix.txt", adj_matrix_txt)

	// task 3: create first ranking
	rank1_prestige_values := []int{}
	ranking1 := make(map[string]int)
	for i := 0; i < len(adj_matrix[0]); i++ {
		sum := 0
		for j := 0; j < len(adj_matrix); j++ {
			site_prestige := 1
			if strings.Contains(site_list[j], "game") || (strings.Contains(site_list[j], "evil-west")) {
				site_prestige = 3
			} else if strings.Contains(site_list[j], "career") || (strings.Contains(site_list[j], "contact")) {
				site_prestige = 2
			}
			// don't calculate cycles
			if i != j {
				sum += adj_matrix[j][i] * site_prestige
			}
		}
		rank1_prestige_values = append(rank1_prestige_values, sum)
		ranking1[site_list[i]] = sum
	}

	fmt.Println("ranking 1:")
	ranking1 = sort_map_by_value(ranking1)

	save_json_file("ranking_1.json", ranking1)

	// save updated labels with prestige
	prestige_labels_txt := ""
	for i := 0; i < len(site_list); i++ {
		prestige_labels_txt += fmt.Sprintf(",%d [%d]", i, rank1_prestige_values[i])
	}

	save_file("ranking1_prestige_labels.txt", prestige_labels_txt)

	// task 4-5: create second ranking
	rank2_prestige_values := []float64{}
	ranking2 := make(map[string]float64)
	for i := 0; i < len(adj_matrix[0]); i++ {
		sum := 0.0
		for j := 0; j < len(adj_matrix); j++ {
			site_prestige := 1.0
			if strings.Contains(site_list[j], "game") || (strings.Contains(site_list[j], "evil-west")) {
				site_prestige = 3
			} else if strings.Contains(site_list[j], "career") || (strings.Contains(site_list[j], "contact")) {
				site_prestige = 2
			}
			// don't calculate cycles
			if i != j {
				Nv := 0.0
				for k := 0; k < len(adj_matrix[j]); k++ {
					Nv += float64(adj_matrix[j][k])
				}
				if Nv != 0 {
					sum += (1.0 / Nv) * site_prestige
				}
			}
		}
		rank2_prestige_values = append(rank2_prestige_values, sum)
		ranking2[site_list[i]] = sum
	}

	fmt.Println("ranking 2:")
	ranking2 = sort_map_by_float_value(ranking2)

	save_json_file("ranking_2.json", ranking2)

	// save updated labels with prestige
	prestige_labels_txt = ""
	for i := 0; i < len(site_list); i++ {
		prestige_labels_txt += fmt.Sprintf(",%d [%.2f]", i, rank2_prestige_values[i])
	}

	save_file("ranking2_prestige_labels.txt", prestige_labels_txt)
}
