# KALK
 Make simple online calculators based on YAML files

- Create a directory where you drop kalk.exe in
- In this directory, create a directory "default"
- In "default" create an index.yaml file
- Run kalk.exe to create the calculate

# Status of the Project
Just a fast prototype for now (although functional)

# Basic syntax of yaml
The YAML can have multiple calculaters (split with ---)
Every calculator has a general, input and output section

Output types can be 
- number 
- decimal
- round1
- round2
- ...
- round5
- ceil
- floor

number is an integer
decimal is a float (default)
roundx define the precision point e.g 1,545435 -> round2 -> 1,55
ceil and floor are typical javascript implementations

```
---
general:
  id: eurodollar
  name: Euro To Dollar
input: 
  - name: euro
    humanName: Euro €
    default: 10
  - name: exchangeRate
    humanName: Exchange Rate
    default: 1.22
output:
  - name: dollar
    humanName: Dollar $
    formula: euro*exchangeRate
---
```

# License
MIT, refer to the license file
Code is not pretty because it's more a script then a full fledged tool

# Why
## What the hell is KALK
Kalk is the dutch word for lime, the ingredient in chalk (for writing on blackboards). It's also conveniently the letter c=k in the word calc

## Why YAML
Not to reinvent the wheel and to have easy parse and extension capabilities