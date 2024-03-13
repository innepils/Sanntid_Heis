##### Error-handling
# Cost-function / Assigning new order
-> If an error happens istead of assingning an order, the node should restart, i.e. the primary should quit, the backup should take over and become the new backup. Then a new backup i started. (PRROCESS-PAIRS)
    -> This ensures that cab-orders are not lost, but taken over by the backup.

# Network
-> If message is not acknowledge by nodes (who are "known" to be alive) it should try again (a number of times/an amount of time, which should lead to a time-out)
    -> If timed out, the node should be removed from the alive-list.

# Power Loss
-> The elevators current-order and cab-button-states should at "all" times be saved to a local file. This ensures that when the nodes dies (due to power loss), that it can retrieve this information to ensure service guarantee.

# Channels need(?) Thoughts from Linus
-> ch_newLocalOrders, for sending orders from the assigner/cost to the local elev. Old orders should be cleared.
-> 



# Adrian 11.03
- Hva om vi bruker id fra main som en variabel i elevatorstate, slik at 