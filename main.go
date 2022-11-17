package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

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

func read_json_file(name string) []byte {
	jsonFile, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func read_site_map(name string) map[string][]string {
	byteValue := read_json_file(name)
	var result map[string][]string
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func read_site_map_order(name string) map[int]string {
	byteValue := read_json_file(name)

	var result map[int]string
	json.Unmarshal([]byte(byteValue), &result)
	return result
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

func calculate_ranking(adjacency_matix [][]float64, site_list []string) (rank_prestige_values []float64, ranking map[string]float64) {
	visited_sites := make(map[string][]float64)
	for i := 0; i < len(site_list); i++ {
		// first value in array is 1 when site was visited and 0 if not, second value is prestige
		visited_sites[site_list[i]] = []float64{0.0, 0.0}
	}

	site_to_rank := []int{0}
	visited_sites[site_list[0]] = []float64{1.0, 1.0}
	for len(site_to_rank) > 0 {
		current_site_index := site_to_rank[0]
		current_row := adjacency_matix[current_site_index]
		current_prestige := visited_sites[site_list[current_site_index]][1]

		visited_sites[site_list[current_site_index]][0] = 1
		for j := 0; j < len(current_row); j++ {
			// not calulate cycles
			calculated_site_index := site_list[j]
			if visited_sites[calculated_site_index][0] == 0 {
				visited_sites[calculated_site_index][1] += current_row[j] * current_prestige
				if current_row[j] != 0 {
					site_to_rank = append(site_to_rank, j)
				}
			}
		}
		site_to_rank = site_to_rank[1:]
	}

	rank_prestige_values = []float64{}
	ranking = make(map[string]float64)

	for i := 0; i < len(site_list); i++ {
		key := site_list[i]
		ranking[key] = visited_sites[key][1]
		rank_prestige_values = append(rank_prestige_values, visited_sites[key][1])
	}
	return
}

func main() {
	source_site_name := "https://flyingwildhog.com/careers/"
	site_map, in_order_site_map := collect_data(source_site_name, "flyingwildhog.com")
	// task 1: print site map - for first run
	save_json_file("site_list.json", site_map)
	save_json_file("site_list_order.json", in_order_site_map)

	// task 1: read existing site structure - for next runs
	// site_map := read_site_map("site_list.json")
	// in_order_site_map := read_site_map_order("site_list_order.json")
	// task 1: create site list
	site_list := []string{}
	// for i := 0; i < len(in_order_site_map); i++ {
	for i := 0; i < len(in_order_site_map); i++ {
		site_name := in_order_site_map[i]
		site_list = append(site_list, site_name)

	}

	out_string := ""
	for _, data := range site_list {
		out_string += data + "\n"
	}

	save_file("site_list.txt", out_string)

	// task 2: create adjacency matrix for site_map
	adj_matrix_txt := ""
	adj_matrix := [][]float64{}
	for i := 0; i < len(site_list); i++ {
		adj_matrix_txt += fmt.Sprintf(",%d", i)
	}
	adj_matrix_txt += "\n"
	for i := 0; i < len(site_list); i++ {
		column := site_list[i]
		adj_matrix_txt += fmt.Sprintf("%d", i)
		column_values := []float64{}
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
	rank1_prestige_values, ranking1 := calculate_ranking(adj_matrix, site_list)

	fmt.Println("ranking 1:")
	ranking1 = sort_map_by_float_value(ranking1)

	save_json_file("ranking_1.json", ranking1)

	// save updated labels with prestige
	prestige_labels_txt := ""
	for i := 0; i < len(site_list); i++ {
		prestige_labels_txt += fmt.Sprintf(",%d [%f]", i, rank1_prestige_values[i])
	}

	save_file("ranking1_prestige_labels.txt", prestige_labels_txt)

	// task 4-5: create second ranking

	updated_adj_matrix := [][]float64{}

	for i := 0; i < len(site_list); i++ {
		column := site_list[i]
		column_values := []float64{}
		for j := 0; j < len(site_list); j++ {
			row := site_list[j]
			if slices.Contains(site_map[column], row) {
				Nv := 0.0
				for k := 0; k < len(adj_matrix[i]); k++ {
					Nv += float64(adj_matrix[i][k])
				}
				if Nv != 0.0 {
					column_values = append(column_values, (1.0 / Nv))
				} else {
					column_values = append(column_values, 0.0)
				}
			} else {
				column_values = append(column_values, 0.0)
			}
		}
		updated_adj_matrix = append(updated_adj_matrix, column_values)
	}

	rank2_prestige_values, ranking2 := calculate_ranking(updated_adj_matrix, site_list)

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
