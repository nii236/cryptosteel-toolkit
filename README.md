# Cryptosteel Generator

The Cryptosteel wallet is an amazing cold storage solution. There is, however, no easy way to convert from a truncated mnemonic to a full one.

Until now!

There are three main functions of this project.

## Generate a mnemonic
```
./cryptosteel-generator --gen

flag sure ramp inquiry panel jaguar balcony merit tunnel profit obtain public comfort dinosaur throw dignity explain symptom cart alert much wash era sample
```

## Get first derived child key (truncated mnemonic)
```
./cryptosteel-generator --t --m "flag sure ramp inqu pane jagu balc meri tunn prof obta publ comf dino thro dign expla symp cart aler much wash era samp"

Using full mnemonic...
mnemonic: flag sure ramp inquiry panel jaguar balcony merit tunnel profit obtain public comfort dinosaur throw dignity explain symptom cart alert much wash era sample
Extended Key: xprvA492WbcgeJiTHFPe7c3J88H6AfJoRUDCWZ6qEnwTrt7E63idxu2egLYmueSvBsqytP2DZW16975RyG82rrLmCAR17NrnqSFxv6nGvueo3fa
scriptPubKey: 18DLFw44Q2oXpW4adCCZoow7MtR7v3znAR
WIF: L1Ghydhy9CEC9wbm957FHw3SdhmsxB5DAMXmSKWvP4chCvuPWLJp
```

## Get first derived child key (untruncated mnemonic)
```
./cryptosteel-generator --m "flag sure ramp inquiry panel jaguar balcony merit tunnel profit obtain public comfort dinosaur throw dignity explain symptom cart alert much wash era sample"

Using full mnemonic...
mnemonic: flag sure ramp inquiry panel jaguar balcony merit tunnel profit obtain public comfort dinosaur throw dignity explain symptom cart alert much wash era sample
Extended Key: xprvA492WbcgeJiTHFPe7c3J88H6AfJoRUDCWZ6qEnwTrt7E63idxu2egLYmueSvBsqytP2DZW16975RyG82rrLmCAR17NrnqSFxv6nGvueo3fa
scriptPubKey: 18DLFw44Q2oXpW4adCCZoow7MtR7v3znAR
WIF: L1Ghydhy9CEC9wbm957FHw3SdhmsxB5DAMXmSKWvP4chCvuPWLJp
```
