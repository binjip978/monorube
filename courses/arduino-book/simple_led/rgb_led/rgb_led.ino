const int BLED = 9;
const int GLED = 10;
const int RLED = 11;
const int BUTTON = 2;

boolean lastButton = LOW;
boolean currentButton = LOW;
boolean ledMode = 0;


void setup() {
  pinMode(BLED, OUTPUT);
  pinMode(GLED, OUTPUT);
  pinMode(RLED, OUTPUT);
  pinMode(BUTTON, INPUT);
}

void loop() {
//  currentButton = debounce(lastButton);
//  if (lastButton == LOW && currentButton == HIGH) {
//    ledMode++;
//  }
//
//  lastButton = currentButton;
//  if (ledMode == 8) {
//    ledMode = 0;
//  }
  setMode(3);
}

boolean debounce(boolean last) {
  boolean current = digitalRead(BUTTON);
  if (last != current) {
    delay(5);
    current = digitalRead(BUTTON);
  }

  return current;
}

void setMode(int mode) {
  if (mode == 1) {
    digitalWrite(RLED, LOW);
    digitalWrite(GLED, HIGH);
    digitalWrite(BLED, HIGH);
  } else if (mode == 2) {
    digitalWrite(RLED, HIGH);
    digitalWrite(GLED, LOW);
    digitalWrite(BLED, HIGH);
  } else if (mode == 3) {
    digitalWrite(RLED, HIGH);
    digitalWrite(GLED, HIGH);
    digitalWrite(BLED, LOW);
  } else if (mode == 4) {
    analogWrite(RLED, 127);
    analogWrite(GLED, 255);
    analogWrite(BLED, 127);
  } else if (mode == 5) {
    analogWrite(RLED, 255);
    analogWrite(GLED, 127);
    analogWrite(BLED, 127);
  } else if (mode == 6) {
    analogWrite(RLED, 127);
    analogWrite(GLED, 127);
    analogWrite(BLED, 255);
  } else if (mode == 7) {
    analogWrite(RLED, 170);
    analogWrite(GLED, 170);
    analogWrite(BLED, 170);
  } else {
    digitalWrite(RLED, LOW);
    digitalWrite(GLED, LOW);
    digitalWrite(BLED, LOW);
  }
}
