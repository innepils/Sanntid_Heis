# Todo
### Adrian:

### Linus:
- die?

### Simon:

##### Kodekvalitet, hva som må endres
# Assigner
- [x] Endre if-statements til switch
- [x] Endre allOrders til allRequests
# Backup
- [ ] Change comments. They should describe more than the code. "Open file for reading" befor e a file-read gives no new infomration.
# Config
- [x] Deleted CV
# Cost
- [x] changed var-name
# Elevator
- Changed name on MAP (elevState) etc.
# Elevator_io
- [ ]
# FSM
- [ ] print Elevator when it brings new information
- [ ] remove "pair" and use pointers. THIS IS IMPORTANT
- [x] When switch-case only has one case, should we use if instead? Might be more clear that it is a State-change, if it is a switch-case? NEI tror jeg. Tydeligere med switch
# NETWORK packages
- [ ] Make some comments explaining the modules (if needed? it is given code.)
# Requests
- [ ] change according to FSM. Use pass-by-reference and not value.
# Main
- [ ] Heartbeat- and networkstuff should be in a function in a goroutine (looks messy in main atm).
 

- localElevator sendes mye inn i funksjoner (spesielt i requests), bør bruke peker og referanse istedenfor kopi. (hastighet og plass)

##### Project description, hva har vi gjort og hva mangler #####
# The button lights are a service guarantee
  - [] Mangler logikk for at hvis en heis ikke fullfører en ordre (hall call button) på gitt tid (e.g. 30sek?), så må andre heiser ta over. Skal watchdog implementeres f.eks.?
  - [] Cab calls fungerer som spesifisert (men vi må sikre at backup og lagring i fil fungerer helt.)

# No calls are lost
  - [] Test forskjellige failures som: losing network entirely, software-crash (watchdog?), obstruction, tap av strøm til både motor og hele noden.
    - [ ] Ved restart (etter crash) hentes cab-orders inn igjen.
    - [x] Når en node er alene på nettverket, skal den fortsatte å fullføre ordre, samt ta nye.
    - [ ] Noden skal IKKE måtte restartes manuelt. (Bør implementere restart i software ved ingen ordre og alene på nettverket.)

# The light and buttons should function as expected
  - [x] Hall call button henter en heis.
  - [ ] Hall call button lights skal være lik på alle nodene (når man er på nettverket med andre noder. Evt   packet loss skal bare føre til noe forsinkelse)
  - [x] Cab button lights skal ikke deles 
  - [x] Knappelys skrus på så fort som mulig (lov å anta at en kunde kan trykke igejn ved ingen lys.)
  - [x] Knappelys skrus av når ordren er fullført.

  # The door should function as expected
  - [x] Lyset simulerer åpen dør (3sek), og skal IKKE gå på når heisen beveger seg.
  - [x] Obstruction hindrer døra i å lukkes, og har ingen påvirkning når heisen beveger seg.
 
  # An individual elevator should behave sensibly and efficiently
  - [x] Heisen skal IKKE stoppe overalt bare for sikkerhetsskyld.
  - [x] Hall call buttons lights som skrus av skal bety at kunder er hentet, og kundene går BARE på om heisen går i retningen de skal.
  - [] Hvis heisen går i en retning den ikke skal (fordi en kunde endret mening og retning, samt det ikke eksisterer andre ordre i den orginale retningen), så skal heisen "ANNONESERE" retningsendring og holde døre åpen i nye 3 sekunder. 

  ### Secondary requirements ###
  # Calls should be served as efficiently as possible.
  - [x] Implementere en fungerende cost-funksjon.
  - [ ] Gjøre koden mer effektiv, bruke pekere, ikke kopiere structs om unødvendig, osv.


  ##### LOVLIGE ANTAGELSER #####
  - Det er alltid MINST EN heis som ikke er påvirket av failure (inkludert obstruction).
  - Det testes med 3 heiser, ikke flere.

