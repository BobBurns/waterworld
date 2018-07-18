char dataString[50] = {0};
int val = 0;
int analogIn = 0;
int dout = 2;
int in_byte = 0;
int i;


void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  pinMode(dout, OUTPUT);
  digitalWrite(dout, LOW);

}

void loop() {
  // put your main code here, to run repeatedly:
  val = analogRead(analogIn);
  
  sprintf(dataString,"%d",val); 
  Serial.println(dataString);   // send the data
  // psuedo delay 
  for (i = 0; i < 10000; i++) {
    if (Serial.available() > 0) {
      in_byte = Serial.read();
      if (in_byte == 0xb0) {
        digitalWrite(dout, HIGH);
      } else if (in_byte == 0xb1) {
        digitalWrite(dout, LOW);
      }
    }
  }

}
