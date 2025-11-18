/*
Package shell provides a builder-based approach to creating and executing shell commands, extending the functionality
provided by os/exec. This functionality is provided via the CommandGenerator interface (instantiated via
NewCommandGenerator()).

If you do not need to modify commands between their creation and execution time, you may wish to use the
internal/console package instead.
*/
package shell
