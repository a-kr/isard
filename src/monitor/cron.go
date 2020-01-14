package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type CronRule struct {
	Minutes []int
	Hours   []int
}

func IntarrayContains(a []int, x int) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func (r *CronRule) Matches(t time.Time) bool {
	m := t.Minute()
	h := t.Hour()

	return ((r.Minutes == nil || IntarrayContains(r.Minutes, m)) &&
		(r.Hours == nil || IntarrayContains(r.Hours, h)))
}

func ParseCronRule(rule string) (*CronRule, error) {
	var err error
	parts := strings.Fields(rule)
	if len(parts) != 2 {
		return nil, fmt.Errorf("rule must have form '<minutes> <hours>': %s", rule)
	}

	parsePart := func(rulePart string, maxValue int) ([]int, error) {
		if rulePart == "*" {
			return nil, nil
		}
		i, err := strconv.Atoi(rulePart)
		if err == nil {
			return []int{i}, nil
		}
		if strings.HasPrefix(rulePart, "*/") {
			modulo, err := strconv.Atoi(rulePart[2:])
			if err != nil {
				return nil, err
			}
			result := []int{}
			for i := 0; i < maxValue; i++ {
				if i%modulo == 0 {
					result = append(result, i)
				}
			}
			return result, nil
		}
		return nil, fmt.Errorf("Invalid cron rule part: %s", rulePart)
	}
	r := &CronRule{}

	r.Minutes, err = parsePart(parts[0], 60)
	if err != nil {
		return nil, err
	}
	r.Hours, err = parsePart(parts[1], 24)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func MustParseCronRule(rule string) *CronRule {
	r, err := ParseCronRule(rule)
	if err != nil {
		panic(err)
	}
	return r
}

func Cron() {
	lastMinuteNumber := time.Now().Minute()

	lastItemsWithCronRules := []string{}

	for now := range time.Tick(1 * time.Second) {
		currentMinuteNumber := now.Minute()
		if lastMinuteNumber == currentMinuteNumber {
			continue
		}
		lastMinuteNumber = currentMinuteNumber

		items, err := FindItems()
		if err != nil {
			log.Printf("error: cannot find items: %s", err)
			continue
		}

		itemsWithCronRules := []string{}

		for _, item := range items {
			if item.CronRule == "" {
				continue
			}
			rule, err := ParseCronRule(item.CronRule)
			if err != nil {
				log.Printf("error: cannot parse cron rule for %s: %s", item.Name, err)
				continue
			}
			itemsWithCronRules = append(itemsWithCronRules, item.Name)
			if rule.Matches(now) {
				go func(itemName string) {
					defer func() {
						if err := recover(); err != nil {
							log.Printf("[%s] panicked: %s", itemName, err)
						}
					}()
					err := Collect(itemName)
					if err != nil {
						log.Printf("[%s] collect error: %s", itemName, err)
					}
				}(item.Name)
			}
		}

		if !SameStringArrays(lastItemsWithCronRules, itemsWithCronRules) {
			lastItemsWithCronRules = itemsWithCronRules
			log.Printf("Items with cron rules: %v", itemsWithCronRules)
		}
	}
}
