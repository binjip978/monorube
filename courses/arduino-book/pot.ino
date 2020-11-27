const int TEMP = A0;
int val = 0;

void setup() {
  Serial.begin(9600);
}

void loop() {
  val = analogRead(TEMP);
  double temp = (double)val / 1024;
  temp = temp * 5;
  temp = temp - 0.5;
  temp = temp * 100;
  Serial.println(temp);
  delay(1000);
}
