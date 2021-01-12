import 'dart:io';


void info(String message) {
  log("INFO", message);
}

void error(String message) {
  log("ERROR", message);
}

void warning(String message) {
  log("WARNING", message);
}

void log(String level, String message) {
  var now = new DateTime.now().toUtc().toIso8601String();
  stdout.writeln("$now [$level] $message\n");
  stdout.flush();
}