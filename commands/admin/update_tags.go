package admin

import (
	"context"
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) updateTags(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) == 0 {
		_, err = ctx.Send("No attachments given.", nil)
		return
	}

	// get the first attachment
	resp, err := http.Get(ctx.Message.Attachments[0].URL)
	if err != nil {
		_, err = ctx.Sendf("I couldn't download the attachment:\n```%v```", err)
		return
	}
	defer resp.Body.Close()

	records, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		_, err = ctx.Sendf("I couldn't parse the attachment as a CSV file:\n```%v```", err)
		return
	}

	// parse into something Usable
	type term struct {
		ID   int
		Tags []string
	}
	toUpdate := []term{}

	displayTags := []string{}

	for i, r := range records {
		if len(r) < 2 {
			_, err = ctx.Sendf("Record %v isn't the correct length, aborting.", i+1)
			return
		}

		id, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			_, err = ctx.Sendf("Couldn't parse the ID in record %v.", i+1)
			return err
		}
		tags := []string{}
		for _, t := range strings.Split(r[1], ",") {
			// add display form to tags
			alreadyInTags := false
			for _, tag := range displayTags {
				if tag == strings.TrimSpace(t) {
					alreadyInTags = true
					break
				}
			}
			if !alreadyInTags {
				displayTags = append(displayTags, strings.TrimSpace(t))
			}

			tags = append(tags, strings.TrimSpace(strings.ToLower(t)))
		}

		toUpdate = append(toUpdate, term{id, tags})
	}

	m, err := ctx.Sendf("About to update the tags for %v terms (totaling %v tags). Continue?", len(toUpdate), len(displayTags))
	if err != nil {
		return err
	}
	yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID)
	if !yes || timeout {
		_, err = ctx.Send("Cancelled.", nil)
		return
	}

	for _, t := range toUpdate {
		_, err = c.DB.Pool.Exec(context.Background(), "update terms set tags = $1 where id = $2", t.Tags, t.ID)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
		c.Sugar.Debugf("Updated %v's tags to %v", t.ID, t.Tags)
	}

	// hehe numbers
	var count int64
	for _, t := range displayTags {
		ct, err := c.DB.Pool.Exec(context.Background(), `insert into public.tags (normalized, display) values ($1, $2)
		on conflict (normalized) do update set display = $2`, strings.ToLower(t), t)
		if err != nil {
			c.Sugar.Errorf("Error adding tag: %v", err)
		}

		count += ct.RowsAffected()
	}

	_, err = ctx.Sendf("Complete! Updated %v terms with %v unique tags.", len(toUpdate), count)
	return
}
