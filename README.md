# Credit Card Spends Tracker (WIP)
A CLI tool to fetch and store your credit card transaction details from your inbox. <b>(India only)</b>

The motivation behind this is to avoid the usage of proprietary apps like CRED or Axio that offer similar functionality but might sell your data. This application is open source and can be built from source if you feel the need to do so. 
An alternative would be an open source mobile application that reads your SMS data which seems to have convuluted permission issues. I made this for myself to quickly check my spends across all my cards.

It does this by making use of the [Gmail API](https://developers.google.com/gmail/api/guides/) to fetch your emails using your <b>own</b> GCP credentials ensuring that you're the <b>only</b> one with access to your emails.

## Features
- Fetch and view your credit card transactions between any date range (currently supports Axis and HDFC cards)
- Displays spend totals by credit card and merchant
- Caches fetched transactions in a local SQLite database to prevent duplicate fetch calls in the future
- Allows you to manually add aliases for merchant names for better organization of your transactions

## Quickstart

1. Clone the repo and run `go build .` to build the application
2. Follow [these](https://developers.google.com/gmail/api/quickstart/go#set_up_your_environment) steps till you get your hands on your `credentials.json` file. Ensure this file is in the root directory.
3. Run `.\credit-card-spends-tracker.exe --fetch 2024-10-10 2024-10-14` (fetches transactions between Oct 10th and Oct 14th)
4. During first time set up, you will have to log in to your respective google account as prompted by the application.

## Commands

1. Fetch transactions between two dates (inclusive)
   
`credit-card-spends-tracker.exe --fetch {YYYY-MM-DD} {YYYY-MM-DD}`

Example: `.\credit-card-spends-tracker.exe --fetch 2024-10-10 2024-10-14`

2. Add an alias for a merchant

`credit-card-spends-tracker.exe --alias {alias} {merchant_name}`

Example: `.\credit-card-spends-tracker.exe --alias "WWW SWIGGY COM" "Swiggy"`

> [!NOTE]  
> An alias is a merchant tag used by a merchant in a transaction email. A merchant can have multiple aliases which requires us to store known aliases for proper categorization.

## Roadmap
1. Add error handling for improper CLI arguments
2. Add unit tests
3. Improve stdout formatting  
