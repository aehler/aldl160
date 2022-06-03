# aldl160
ALDL windows application

Reader and logger of a raw ALDL (GM OBD1) data stream at 160 baud, redirected at any com port by arduino at 115200baud

This application reads serial port in realtime, parses aldl format, decrypts data for 3.1L GM Van (Vin-D).
Can show online sensor data, errors and other info sent by EPROM

Buids with https://developer.fyne.io/index.html Read docs for infrastructure details.

Arduino sketch is as simple as it can be, just echoing a pin, connected to E pin of ALDL port, to serial port at 115200baud

Test vehicle is submitting data at 5V, thus can be directly attached to arduino port.
For more info on arduino schematic and sketches see arduino-aldl repo.

Todo:
- make an arduino sketch to work using interrupts, but need a hardware current filter.
- make viewer of trip logs
- make bluetooth interface
- build app for android
