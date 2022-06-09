streams -> dag (tree?) -> output

Nodes hold last value so they can stop prop if no diff.

Streams - appear to be pushing
* Some might push
* Others might poll and push
* Initial initial value: type default

Expression language for dag: CEL

Add pending: basically a propagating Option.
* Streams to have a timeout, if they don't push within that time, they should push Pending
Expand to be Result - value or an error saying why not.
