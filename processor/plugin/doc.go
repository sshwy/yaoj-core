/*
This package describes how to build a custom processor by creating a shared
library with specific symbol exposed.

Conventionally, plugin is written in go which is compiled to a shared library
(plugin) and is loaded in runtime. However we also plan to support C/C++ plugin
due to its popularity.
*/
package plugin
