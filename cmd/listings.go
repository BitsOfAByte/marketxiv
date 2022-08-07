/*
Copyright © 2022 BitsOfAByte

*/

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BitsOfAByte/marketxiv/backend"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// itemCmd represents the item command
var listingsCmd = &cobra.Command{
	Use:   "listings <server> <item>",
	Short: "Get the current market listings for the specified item",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		hq, _ := cmd.Flags().GetBool("hq")
		limit, _ := cmd.Flags().GetInt("limit")

		serverName := args[0]
		itemName := strings.Join(args[1:], " ")
		searchData := backend.FetchSearch(itemName, "item")

		// Check to see if the item exists
		if len(searchData.Results) == 0 {
			fmt.Println("No results found for " + itemName)
			return
		}

		resultData := searchData.Results[0]
		marketData := backend.FetchMarketItem(serverName, resultData.ID, limit, strconv.FormatBool(hq))

		if len(marketData.Listings) == 0 {
			fmt.Println("No listings found for " + itemName)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Quality", "Price", "Quantity", "Total", "Retainer", "World"})

		// Format and display the data
		for _, listing := range marketData.Listings {
			world := listing.WorldName

			if world == "" {
				world = serverName
			}

			quality := "Normal"
			switch listing.Hq {
			case true:
				quality = "HQ"
			case false:
				quality = "Normal"
			}

			table.Append([]string{
				quality,
				strconv.Itoa(listing.PricePerUnit),
				strconv.Itoa(listing.Quantity),
				strconv.Itoa(listing.Total),
				listing.RetainerName,
				world,
			})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listingsCmd)
	listingsCmd.Flags().Bool("hq", false, "Only fetch high quality listings")
	listingsCmd.Flags().IntP("limit", "l", 5, "Limit the number of listings to show")
}
