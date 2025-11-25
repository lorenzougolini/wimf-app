package templates

import (
	"time"
)

func getCardClass(exp time.Time) string {
	if time.Now().After(exp) {
		return "flex items-center justify-between p-4 border rounded-2xl shadow-sm transition-shadow duration-200 bg-red-100 border-red-200 dark:bg-red-900/30 dark:border-red-800"
	}
	return "flex items-center justify-between p-4 bg-white border border-gray-200 rounded-2xl shadow-sm dark:bg-gray-800 dark:border-gray-700 hover:shadow-md transition-shadow duration-200"
}

func getDateClass(exp time.Time) string {
	hoursRemaining := time.Until(exp).Hours()
	days := hoursRemaining / 24

	baseClass := "mt-3 flex item-center text-sm font-medium"

	if hoursRemaining < 0 {
		return baseClass + " text-red-900 dark:text-red-200 font-bold"
	} else if days <= 2 {
		return baseClass + " text-red-600 dark:text-red-400"
	} else if days <= 5 {
		return baseClass + " text-yellow-600 dark:text-yellow-400"
	} else {
		return baseClass + " text-green-600 dark:text-green-400"
	}
}
