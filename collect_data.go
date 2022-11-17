package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/exp/slices"
)

func collect_data(site_url string, allowed_domains string) (site_map map[string][]string, in_order_site_map map[int]string) {
	site_map = make(map[string][]string)
	in_order_site_map = make(map[int]string)
	visited_site_count := 0
	site_root := ""
	site_queue := []string{site_url}
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

		if !slices.Contains(temp_list, current_link) {
			if strings.Contains(current_link, allowed_domains) {
				site_map[site_root] = append(temp_list, current_link)
			}
		}
		// check if site exist in site_map
		if _, site_exist := site_map[current_link]; !site_exist {
			site_queue = append(site_queue, link)
		}
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		site_root = r.URL.String()
		site_map[site_root] = []string{}
		in_order_site_map[visited_site_count] = site_root
		visited_site_count++
	})

	site_queue = append(site_queue, site_url)
	for len(site_queue) > 0 {
		current_link := site_queue[0]
		site_queue = site_queue[1:]
		c.Visit(current_link)
	}
	fmt.Println("Visited ", visited_site_count, " sites")
	return
}
