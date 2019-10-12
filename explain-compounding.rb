#!/usr/bin/env ruby

require 'humanize'

initial_amount = 1_26_000
rate = 9.0
years = 10

puts "Use defaults - Amount: #{initial_amount}, Rate: #{rate}, Years: #{years}?"
default = gets

if default.include? "n"
  puts "Enter amount: "
  initial_amount = gets.to_i

  puts "Enter rate of interest: "
  rate = gets.to_f

  puts "Enter total years: "
  years = gets.to_i
end

amount = initial_amount

puts "You invested \"#{amount.humanize}\" for #{years} years at #{rate}% compounded"

years.times do |i|
  amount = (amount * ( 1 + (rate / 100.0) )).round(2)
  puts "At end of #{i+1} year: #{amount} or \"#{amount.humanize}\""
end

puts "If you invest it for another #{years} years: #{amount * (( 1 + (rate / 100.0) ) ** years)}"
