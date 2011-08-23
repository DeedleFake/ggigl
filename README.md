ggigl
=====

ggigl (Go Game In Go Lang) is an implementation of the [Go board game][gogame] in the [Go programming language][golang].

Status
------

Game is playable, but can't calculate score. In fact, there's no way to end the game other than quiting.

For more information, see TODO.

Usage
-----

> ggigl [options]

<dl>
	<dt>-size=&lt;num&gt; (Default: 19)</dt>
		<dd>The board size to use. Only accepts 9 or 19.</dd>
	<dt>-handicap=&lt;num&gt; (Default: 0)</dt>
		<dd>The handicap. 0 is no handicap; maximum is 19.</dd>
	<dt>-komi=&lt;num&gt; (Default: -1)</dt>
		<dd>Komi. Anything less than 0 will have different results depending on the handicap settings. If there's no handicap, then komi will be 5.5. If there is a handicap, komi will be 0.</dd>
</dl>

Authors
-------

 * DeedleFake

[gogame]: http://www.wikipedia.com/wiki/Go_(board_game)
[golang]: http://www.golang.org
