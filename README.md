ggigl
=====

ggigl (Go Game In Go Lang) is an implementation of the [Go board game][gogame] in the [Go programming language][golang].

Usage
-----

> ggigl [options]

<dl>
	<dt>-size=<num> (Default: 19)</dt>
		<dd>The board size to use. Only accepts 9 or 19.</dd>
	<dt>-handicap=<num> (Default: 0)</dt>
		<dd>The handicap. 0 is no handicap; maximum is 19.</dd>
	<dt>-komi=<num> (Default: -1)</dt>
		<dd>Komi. Anything less than 0 will have different results depending on the handicap settings. If there's no handicap, then komi will be 5.5. If there is a handicap, komi will be 0.</dd>
</dl>

Authors
-------

 * DeedleFake: <yisszev at beckforce dot com>

[gogame]: http://www.wikipedia.com/wiki/Go_(board_game)
[golang]: http://www.golang.org
