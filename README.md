# Sequence - Counter as a Service
Track and increment all of the things 👨🏿‍💻

Wanna name stuff, but want the names to follow a numbered sequence? tada
![number_sequence](https://user-images.githubusercontent.com/673382/42729381-730a61e2-87a4-11e8-9b4e-da34c3d56005.png)

Uses Dynamodb to increment values in order to have sequential naming schemes 

Our primary use case for this is AWS Auto Scaling group instance names - each instances will name itself after hitting this endpoint and getting a position in the sequence

## Deployment steps
1. Use [Cloudformation template in CF directory](https://github.com/Sjeanpierre/Sequence/blob/master/CF/dynamo_create.json) to create DynamoDB table in chosen region
2. `make`
3. `sls deploy`
4. ????
5. profit
