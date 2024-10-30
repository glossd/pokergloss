# Profile

#### Sign Up Flow

1. MongoDB: checks username uniqueness
2. Creates user in Google Cloud Identity Platform (emailVerified true without actual verification)
3. Saves username in MongoDB
4. Sets username custom claim in GCIP  
5. Sends bonus chips to bank service
6. Gets custom token from GCIP

