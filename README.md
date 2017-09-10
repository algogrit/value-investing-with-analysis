Value Investing and Analysis
----------------------------

This repo contains scraping tools, analysis and visualizations based on formulae and analysis as done in the book "Value Investing and Behavioral Finance" by Parag Parikh.


## Setup

  $ ./scripts/setup.sh


## TODO

  - Scrape data for all the companies listed in NSE
    * Start with scraping financial data from Screener
    * Scrape total stocks and the price per day for last 10 years, if possible along with P/E ratios
  - Start combing the data and inserting into database
    * Create a postgres database
    * Insert data for each company normalized
    * Insert analyzed data like fundamental and speculatory influence over the P/E ratio
  - Visualize the data using simple line charts
