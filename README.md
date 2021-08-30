## Signing and verifying transactions using wallets

The whole point of signing and verifying transactions is to ensure that the person initiating the transaction
is not using funds which do not belong to them.

Let's say there are two people on the network: Andrew and an impostor. Andrew has received 50 GPCoins from a
previous transaction output. Without signatures and verification, the impostor can pretend to be Andrew by
initiating a transaction using Andrew's previous transaction outputs. There is no way to verify whether it is
Andrew initiating the transaction & whether the funds belong to the person initiating the transaction. We need
some sort of security mechanism for ensuring that (a) no one else can use money that I own, and (b) I cannot
use the money of other people.

Enter wallets. Each person has their own **wallet**, which holds a _private key_, as well as an **address**,
which is the _public key_ associated with their wallet/private key. A security system is implemented as follows:

1. When initiating a transaction, I specify transaction inputs, and sign the transaction (i.e., create a signature)
   using my wallet/private key.
2. Each transaction input has a 1-to-1 correspondence with a previous transaction output. Therefore, we can find
   all the previous transaction outputs corresponding to my current inputs.
3. The previous transaction outputs have an address, which is a public key. If this address/public key is indeed mine
   (i.e., if I actually own these funds), I should be able to use this address to verify that my wallet/private key was
   used to create the newly initiated transaction’s signature.
4. Upon successful verification, we can now add this new transaction to the mempool and subsequently the blockchain.

We can now see how such a system would work against an "impostor" who tries to use someone else's funds:

1. The impostor says “I’m Andrew, and I want to initiate this transaction. So I’ll use his previous funds.”
2. The impostor’s wallet/private key is not associated with Andrew's address (public key).
3. When they try to initiate the transaction, they will sign the transaction with their own wallet/private key.
4. However, Andrew's address (public key) will fail to verify the signature, since it was signed using the impostor’s
   wallet/private key. Therefore, we know that the funds don’t belong to this person & it is not actually Andrew
   initiating the transaction, so we block it from being added to the mempool.
