# Todo
### Alle
 - Kommenter all kode man har laget ut ifra kode-kvalitetsperspektiv (gjerne se pensumboka)
    - det viktigste er at kommentaren skal gi ytterligere informasjon, ikke gjengi direkte hva koden sier.
      - Slik kommentering (sistnevnte) er verre enn ingen kommentering.
    - Ikke benytt "inline-" kommentering om det ikke er helt nødvendig, altså kommentarer på slutten av en kode linje (til høyre)
- Er visst dårlig å sende channelse fra funksjon til funksjon. De kanalene som bare brukes mellom assigner og fsm (for eksempel); kan de lages lokalt i fsm istedenfor? Altså mindre channel-passing
- Se over hvilke moduler som importerer andre, direkte import-cycle går ikke an men vi bør unngå en unødvendig link frem og tilbake mellom moduler om det er mulig.

### Adrian:
 - [ ] Change as much pass-by-value to pass-by-reference. E.g. in elevator.go
### Linus:
- die? i cost
  - lage error-funksjoner istedenfor å ha if-statement/3-linjer hver gang
- assigner linje 42, hva er 4?
- assigner linje 79 og for loops generelt, skal vi ha i og j, eller floor og byttontype?
- backup function SaveBackupToFile. Skal vi ha floor og byttontype der det står request og 2?
- config linje 24. "BackupRecieveAddr" (vi endret fra "recieveAddr"), skal den være "DefaultPortBackup"?
- backup linje 129, må den være 2.5 sekunder? Den er 5 sekunder lenger opp på linje 107.
- 


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
- [x] Foor loop for localeelevator is ridicioulous
# Elevator
- [x] Changed name on MAP (elevState) etc.
- [ ] change functionname from ElevButtonToString to ElevButtonTypeToString?? Brukes den i det hele tatt?
# Elevator_io
- [ ] Should short one-letter variable names be renamed? It is given code..
# FSM
- [ ] print Elevator when it brings new information
- [x] remove "pair" and use pointers. THIS IS IMPORTANT
- [x] When switch-case only has one case, should we use if instead? Might be more clear that it is a State-change, if it is a switch-case? NEI tror jeg. Tydeligere med switch
- [ ] There is a lot of duplicated code, due to the state maching being state-driven. Also the stop-button is really far to long. Hear with the group if/what can be made into routines (functions)
    - also, delete unnecessary from stop-button if it works (test when all functionality works)
# Heartbeat
- [ ] See if we can divide into two functions. Should not have a function that does 3 things
# NETWORK packages
- [ ] Make some comments explaining the modules (if needed? it is given code.)
# Requests
- [x] change according to FSM. Use pass-by-reference and not value.
# Main
- [x] Heartbeat- and networkstuff should be in a function in a goroutine (looks messy in main atm).

##### Project description, hva har vi gjort og hva mangler #####
# The button lights are a service guarantee
  - [ ] Mangler logikk for at hvis en heis ikke fullfører en ordre (hall call button) på gitt tid (e.g. 30sek?), så må andre heiser ta over. Skal watchdog implementeres f.eks.?
  - [ ] Cab calls fungerer som spesifisert (men vi må sikre at backup og lagring i fil fungerer helt.)

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
  - [ ] Hvis heisen går i en retning den ikke skal (fordi en kunde endret mening og retning, samt det ikke eksisterer andre ordre i den orginale retningen), så skal heisen "ANNONESERE" retningsendring og holde døre åpen i nye 3 sekunder. 

  ### Secondary requirements ###
  # Calls should be served as efficiently as possible.
  - [x] Implementere en fungerende cost-funksjon.
  - [ ] Gjøre koden mer effektiv, bruke pekere, ikke kopiere structs om unødvendig, osv.


  ##### LOVLIGE ANTAGELSER #####
  - Det er alltid MINST EN heis som ikke er påvirket av failure (inkludert obstruction).
  - Det testes med 3 heiser, ikke flere.

