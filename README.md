Golang SOS Injector
===

This is a basic Program for injecting Air Quality Data automatically obtained from the "Umweltbundesamt" ([www.umweltbundesamt.de](www.umweltbundesamt.de)) into a SOS ([Sensor Observation Service](https://en.wikipedia.org/wiki/Sensor_Observation_Service)) Server.

The Program is used in the Final Assignment in the course "SII - Spatial Information Infrastructures" in WWU Muenster ([link][1]).

**Due to time restrictions the Program is programmed specific to this Use Case and not optimized to be portable in any way.**

There is **NO** Documentation available as this program is not intended to be used outside of the Assignment.
 
 [1]: https://uvlsf.uni-muenster.de/qisserver/rds?state=verpublish&status=init&vmfile=no&publishid=230590&moduleCall=webInfo&publishConfFile=webInfo&publishSubDir=veranstaltung " "

##Dependencies
 - goscv Library ([github repo][https://github.com/gocarina/gocsv])

##TODO:
 - Error Handling
 - Log File Parsing
 - CRON Job for automatic Insertion of Data
 - Improve XML Templates (remove Dummy Variables)
 - Documentation