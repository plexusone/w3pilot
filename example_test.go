package webpilot_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/plexusone/webpilot"
)

func Example() {
	ctx := context.Background()

	// Launch browser
	pilot, err := webpilot.Launch(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	// Navigate to a page
	if err := pilot.Go(ctx, "https://example.com"); err != nil {
		log.Fatal(err)
	}

	// Find and click a link
	link, err := pilot.Find(ctx, "a", nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := link.Click(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Get page title
	title, err := pilot.Title(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Page title:", title)
}

func Example_headless() {
	ctx := context.Background()

	// Launch headless browser
	pilot, err := webpilot.LaunchHeadless(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	// Navigate
	if err := pilot.Go(ctx, "https://example.com"); err != nil {
		log.Fatal(err)
	}

	// Take screenshot
	data, err := pilot.Screenshot(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	if err := os.WriteFile("screenshot.png", data, 0600); err != nil {
		log.Fatal(err)
	}
}

func Example_formInteraction() {
	ctx := context.Background()

	pilot, err := webpilot.Browser.Launch(ctx, &webpilot.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	// Navigate to a form page
	if err := pilot.Go(ctx, "https://example.com/login"); err != nil {
		log.Fatal(err)
	}

	// Fill in username
	username, err := pilot.Find(ctx, "input[name='username']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := username.Type(ctx, "myuser", nil); err != nil {
		log.Fatal(err)
	}

	// Fill in password
	password, err := pilot.Find(ctx, "input[name='password']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := password.Type(ctx, "mypassword", nil); err != nil {
		log.Fatal(err)
	}

	// Click submit
	submit, err := pilot.Find(ctx, "button[type='submit']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := submit.Click(ctx, nil); err != nil {
		log.Fatal(err)
	}
}
