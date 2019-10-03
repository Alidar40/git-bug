package termui

import (
	"bytes"
	"fmt"
	"time"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/util/colors"
	"github.com/MichaelMure/git-bug/util/text"
	"github.com/MichaelMure/gocui"
	"github.com/dustin/go-humanize"
)

const bugTableView = "bugTableView"
const bugTableHeaderView = "bugTableHeaderView"
const bugTableFooterView = "bugTableFooterView"
const bugTableInstructionView = "bugTableInstructionView"

const defaultRemote = "origin"
const defaultQuery = "status:open"

type bugTable struct {
	repo         *cache.RepoCache
	queryStr     string
	query        *cache.Query
	allIds       []entity.Id
	excerpts     []*cache.BugExcerpt
	pageCursor   int
	selectCursor int
}

func newBugTable(c *cache.RepoCache) *bugTable {
	query, err := cache.ParseQuery(defaultQuery)
	if err != nil {
		panic(err)
	}

	return &bugTable{
		repo:         c,
		query:        query,
		queryStr:     defaultQuery,
		pageCursor:   0,
		selectCursor: 0,
	}
}

func (bt *bugTable) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if maxY < 4 {
		// window too small !
		return nil
	}

	v, err := g.SetView(bugTableHeaderView, -1, -1, maxX, 3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
	}

	v.Clear()
	bt.renderHeader(v, maxX)

	v, err = g.SetView(bugTableView, -1, 1, maxX, maxY-3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorWhite
		v.SelFgColor = gocui.ColorBlack

		// restore the cursor
		// window is too small to set the cursor properly, ignoring the error
		_ = v.SetCursor(0, bt.selectCursor)
	}

	_, viewHeight := v.Size()
	err = bt.paginate(viewHeight)
	if err != nil {
		return err
	}

	err = bt.cursorClamp(v)
	if err != nil {
		return err
	}

	v.Clear()
	bt.render(v, maxX)

	v, err = g.SetView(bugTableFooterView, -1, maxY-4, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
	}

	v.Clear()
	bt.renderFooter(v, maxX)

	v, err = g.SetView(bugTableInstructionView, -1, maxY-2, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.BgColor = gocui.ColorBlue

		_, _ = fmt.Fprintf(v, "[q] Quit [s] Search [←↓↑→,hjkl] Navigation [↵] Open bug [n] New bug [i] Pull [o] Push [f] Mark favorite")
	}

	_, err = g.SetCurrentView(bugTableView)
	return err
}

func (bt *bugTable) keybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding(bugTableView, 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	// Down
	if err := g.SetKeybinding(bugTableView, 'j', gocui.ModNone,
		bt.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyArrowDown, gocui.ModNone,
		bt.cursorDown); err != nil {
		return err
	}
	// Up
	if err := g.SetKeybinding(bugTableView, 'k', gocui.ModNone,
		bt.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyArrowUp, gocui.ModNone,
		bt.cursorUp); err != nil {
		return err
	}

	// Previous page
	if err := g.SetKeybinding(bugTableView, 'h', gocui.ModNone,
		bt.previousPage); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyArrowLeft, gocui.ModNone,
		bt.previousPage); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyPgup, gocui.ModNone,
		bt.previousPage); err != nil {
		return err
	}
	// Next page
	if err := g.SetKeybinding(bugTableView, 'l', gocui.ModNone,
		bt.nextPage); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyArrowRight, gocui.ModNone,
		bt.nextPage); err != nil {
		return err
	}
	if err := g.SetKeybinding(bugTableView, gocui.KeyPgdn, gocui.ModNone,
		bt.nextPage); err != nil {
		return err
	}

	// New bug
	if err := g.SetKeybinding(bugTableView, 'n', gocui.ModNone,
		bt.newBug); err != nil {
		return err
	}

	// Open bug
	if err := g.SetKeybinding(bugTableView, gocui.KeyEnter, gocui.ModNone,
		bt.openBug); err != nil {
		return err
	}

	// Pull
	if err := g.SetKeybinding(bugTableView, 'i', gocui.ModNone,
		bt.pull); err != nil {
		return err
	}

	// Push
	if err := g.SetKeybinding(bugTableView, 'o', gocui.ModNone,
		bt.push); err != nil {
		return err
	}

	// Query
	if err := g.SetKeybinding(bugTableView, 's', gocui.ModNone,
		bt.changeQuery); err != nil {
		return err
	}

	// Mark favorite
	if err := g.SetKeybinding(bugTableView, 'f', gocui.ModNone,
		bt.markFavorite); err != nil {
		return err
	}

	return nil
}

func (bt *bugTable) disable(g *gocui.Gui) error {
	if err := g.DeleteView(bugTableView); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err := g.DeleteView(bugTableHeaderView); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err := g.DeleteView(bugTableFooterView); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err := g.DeleteView(bugTableInstructionView); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func (bt *bugTable) paginate(max int) error {
	bt.allIds = bt.repo.QueryBugs(bt.query)

	return bt.doPaginate(max)
}

func (bt *bugTable) doPaginate(max int) error {
	// clamp the cursor
	bt.pageCursor = maxInt(bt.pageCursor, 0)
	bt.pageCursor = minInt(bt.pageCursor, len(bt.allIds))

	nb := minInt(len(bt.allIds)-bt.pageCursor, max)

	if nb < 0 {
		bt.excerpts = []*cache.BugExcerpt{}
		return nil
	}

	// slice the data
	ids := bt.allIds[bt.pageCursor : bt.pageCursor+nb]

	bt.excerpts = make([]*cache.BugExcerpt, len(ids))

	for i, id := range ids {
		excerpt, err := bt.repo.ResolveBugExcerpt(id)
		if err != nil {
			return err
		}

		bt.excerpts[i] = excerpt
	}

	return nil
}

func (bt *bugTable) getTableLength() int {
	return len(bt.excerpts)
}

func (bt *bugTable) getColumnWidths(maxX int) map[string]int {
	m := make(map[string]int)
	m["id"] = 9
	m["status"] = 7

	left := maxX - 5 - m["id"] - m["status"]

	m["summary"] = 10
	left -= m["summary"]
	m["lastEdit"] = 19
	left -= m["lastEdit"]

	m["author"] = minInt(maxInt(left/3, 15), 10+left/8)
	m["title"] = maxInt(left-m["author"], 10)

	return m
}

func (bt *bugTable) render(v *gocui.View, maxX int) {
	columnWidths := bt.getColumnWidths(maxX)

	for _, excerpt := range bt.excerpts {
		summaryTxt := fmt.Sprintf("C:%-2d L:%-2d",
			excerpt.LenComments,
			len(excerpt.Labels),
		)

		var authorDisplayName string
		if excerpt.AuthorId != "" {
			author, err := bt.repo.ResolveIdentityExcerpt(excerpt.AuthorId)
			if err != nil {
				panic(err)
			}
			authorDisplayName = author.DisplayName()
		} else {
			authorDisplayName = excerpt.LegacyAuthor.DisplayName()
		}

		lastEditTime := time.Unix(excerpt.EditUnixTime, 0)

		id := text.LeftPadMaxLine(excerpt.Id.Human(), columnWidths["id"], 1)
		status := text.LeftPadMaxLine(excerpt.Status.String(), columnWidths["status"], 1)
		title := text.LeftPadMaxLine(excerpt.Title, columnWidths["title"], 1)
		author := text.LeftPadMaxLine(authorDisplayName, columnWidths["author"], 1)
		summary := text.LeftPadMaxLine(summaryTxt, columnWidths["summary"], 1)
		lastEdit := text.LeftPadMaxLine(humanize.Time(lastEditTime), columnWidths["lastEdit"], 1)

		if excerpt.IsFavorite {
			_, _ = fmt.Fprintf(v, "%s %s %s %s %s %s\n",
				colors.BlueBg(colors.Cyan(id)),
				colors.BlueBg(colors.Yellow(status)),
				colors.BlueBg(title),
				colors.BlueBg(colors.Magenta(author)),
				colors.BlueBg(summary),
				colors.BlueBg(lastEdit),
			)
		} else {
			_, _ = fmt.Fprintf(v, "%s %s %s %s %s %s\n",
				colors.Cyan(id),
				colors.Yellow(status),
				title,
				colors.Magenta(author),
				summary,
				lastEdit,
			)
		}
	}
}

func (bt *bugTable) renderHeader(v *gocui.View, maxX int) {
	columnWidths := bt.getColumnWidths(maxX)

	id := text.LeftPadMaxLine("ID", columnWidths["id"], 1)
	status := text.LeftPadMaxLine("STATUS", columnWidths["status"], 1)
	title := text.LeftPadMaxLine("TITLE", columnWidths["title"], 1)
	author := text.LeftPadMaxLine("AUTHOR", columnWidths["author"], 1)
	summary := text.LeftPadMaxLine("SUMMARY", columnWidths["summary"], 1)
	lastEdit := text.LeftPadMaxLine("LAST EDIT", columnWidths["lastEdit"], 1)

	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%s %s %s %s %s %s\n", id, status, title, author, summary, lastEdit)

}

func (bt *bugTable) renderFooter(v *gocui.View, maxX int) {
	_, _ = fmt.Fprintf(v, " \nShowing %d of %d bugs", len(bt.excerpts), len(bt.allIds))
}

func (bt *bugTable) cursorDown(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()

	// If we are at the bottom of the page, switch to the next one.
	if y+1 > bt.getTableLength()-1 {
		_, max := v.Size()

		if bt.pageCursor+max >= len(bt.allIds) {
			return nil
		}

		bt.pageCursor += max
		bt.selectCursor = 0
		_ = v.SetCursor(0, bt.selectCursor)

		return bt.doPaginate(max)
	}

	y = minInt(y+1, bt.getTableLength()-1)
	// window is too small to set the cursor properly, ignoring the error
	_ = v.SetCursor(0, y)
	bt.selectCursor = y

	return nil
}

func (bt *bugTable) cursorUp(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()

	// If we are at the top of the page, switch to the previous one.
	if y-1 < 0 {
		_, max := v.Size()

		if bt.pageCursor == 0 {
			return nil
		}

		bt.pageCursor = maxInt(0, bt.pageCursor-max)
		bt.selectCursor = max - 1
		_ = v.SetCursor(0, bt.selectCursor)

		return bt.doPaginate(max)
	}

	y = maxInt(y-1, 0)
	// window is too small to set the cursor properly, ignoring the error
	_ = v.SetCursor(0, y)
	bt.selectCursor = y

	return nil
}

func (bt *bugTable) cursorClamp(v *gocui.View) error {
	_, y := v.Cursor()

	y = minInt(y, bt.getTableLength()-1)
	y = maxInt(y, 0)

	// window is too small to set the cursor properly, ignoring the error
	_ = v.SetCursor(0, y)
	bt.selectCursor = y

	return nil
}

func (bt *bugTable) nextPage(g *gocui.Gui, v *gocui.View) error {
	_, max := v.Size()

	if bt.pageCursor+max >= len(bt.allIds) {
		return nil
	}

	bt.pageCursor += max

	return bt.doPaginate(max)
}

func (bt *bugTable) previousPage(g *gocui.Gui, v *gocui.View) error {
	_, max := v.Size()

	if bt.pageCursor == 0 {
		return nil
	}

	bt.pageCursor = maxInt(0, bt.pageCursor-max)

	return bt.doPaginate(max)
}

func (bt *bugTable) newBug(g *gocui.Gui, v *gocui.View) error {
	return newBugWithEditor(bt.repo)
}

func (bt *bugTable) openBug(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	id := bt.excerpts[y].Id
	b, err := bt.repo.ResolveBug(id)
	if err != nil {
		return err
	}
	ui.showBug.SetBug(b)
	return ui.activateWindow(ui.showBug)
}

func (bt *bugTable) markFavorite(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	id := bt.excerpts[y].Id
	err := bt.repo.MarkFavorite(id)
	if err != nil {
		return err
	}

	return nil
}

func (bt *bugTable) pull(g *gocui.Gui, v *gocui.View) error {
	ui.msgPopup.Activate("Pull from remote "+defaultRemote, "...")

	go func() {
		stdout, err := bt.repo.Fetch(defaultRemote)

		if err != nil {
			g.Update(func(gui *gocui.Gui) error {
				ui.msgPopup.Activate(msgPopupErrorTitle, err.Error())
				return nil
			})
		} else {
			g.Update(func(gui *gocui.Gui) error {
				ui.msgPopup.UpdateMessage(stdout)
				return nil
			})
		}

		var buffer bytes.Buffer
		beginLine := ""

		for result := range bt.repo.MergeAll(defaultRemote) {
			if result.Status == entity.MergeStatusNothing {
				continue
			}

			if result.Err != nil {
				g.Update(func(gui *gocui.Gui) error {
					ui.msgPopup.Activate(msgPopupErrorTitle, err.Error())
					return nil
				})
			} else {
				_, _ = fmt.Fprintf(&buffer, "%s%s: %s",
					beginLine, colors.Cyan(result.Entity.Id().Human()), result,
				)

				beginLine = "\n"

				g.Update(func(gui *gocui.Gui) error {
					ui.msgPopup.UpdateMessage(buffer.String())
					return nil
				})
			}
		}

		_, _ = fmt.Fprintf(&buffer, "%sdone", beginLine)

		g.Update(func(gui *gocui.Gui) error {
			ui.msgPopup.UpdateMessage(buffer.String())
			return nil
		})

	}()

	return nil
}

func (bt *bugTable) push(g *gocui.Gui, v *gocui.View) error {
	ui.msgPopup.Activate("Push to remote "+defaultRemote, "...")

	go func() {
		// TODO: make the remote configurable
		stdout, err := bt.repo.Push(defaultRemote)

		if err != nil {
			g.Update(func(gui *gocui.Gui) error {
				ui.msgPopup.Activate(msgPopupErrorTitle, err.Error())
				return nil
			})
		} else {
			g.Update(func(gui *gocui.Gui) error {
				ui.msgPopup.UpdateMessage(stdout)
				return nil
			})
		}
	}()

	return nil
}

func (bt *bugTable) changeQuery(g *gocui.Gui, v *gocui.View) error {
	return editQueryWithEditor(bt)
}
