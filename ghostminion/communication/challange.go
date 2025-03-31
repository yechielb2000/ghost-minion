package communication

/*
before accepting any communication between the backdoor and the C2,
the backdoor have to answer a challenge to make sure they are who they say they are.
*/

/*
ideas for challenges
 - send special cookie to special route
   maybe something like /instance/<targetID>/key = the cookie (should be delivered at execution time)
*/

func AnswerChallenge() bool {
	return true
}
