package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Akihira77/gojobber/services/5-gig/handler"
	"github.com/Akihira77/gojobber/services/5-gig/types"
	"github.com/Akihira77/gojobber/services/5-gig/util"
	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, _ := NewStore()
	db.Debug().Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	// db.Migrator().DropTable(
	// 	&types.Gig{},
	// )
	// db.AutoMigrate(
	// 	&types.Gig{},
	// )
	// err = types.ApplyDBSetup(db)
	// seedingGig(db)
	if err != nil {
		log.Fatalf("Error applying DB setup:\n%+v", err)
	}

	cld := util.NewCloudinary()
	ccs := handler.NewGRPCClients()
	ccs.AddClient(types.USER_SERVICE, os.Getenv("USER_GRPC_PORT"))

	NewHttpServer(db, cld, ccs)
}

func seedingGig(db *gorm.DB) {
	// Generate and insert dummy data
	titles := []string{
		"Custom T-Shirt Design",
		"Email Marketing Campaign Setup",
		"3D Product Rendering",
		"Creative Copywriting",
		"Website Speed Optimization",
		"Product Photography",
		"Digital Business Card Creation",
		"AI-Powered Chatbot Development",
		"Podcast Editing",
		"Resume and Cover Letter Writing",
		"Interactive Web Banner Design",
		"Infographic Design",
		"Online Course Creation",
		"Data Entry and Virtual Assistance",
		"LinkedIn Profile Optimization",
		"Animated Explainer Videos",
		"Photo Retouching and Enhancement",
		"Social Media Account Management",
		"Custom WordPress Plugin Development",
		"Advanced Excel Data Analysis",
	}
	descriptions := []string{
		"I will create a unique and eye-catching T-shirt design that stands out. Whether it's for personal use or a business, my designs are tailored to meet your specific needs.",
		"Boost your business with a professionally designed email marketing campaign. From template design to automation setup, I ensure your emails reach the right audience at the right time.",
		"Bring your products to life with high-quality 3D rendering. Perfect for online stores, marketing materials, and product presentations. I deliver realistic and detailed renderings that showcase your products in the best light.",
		"Engage your audience with compelling and creative copywriting. I specialize in crafting content that not only captures attention but also drives conversions. Ideal for websites, ads, and social media.",
		"Improve your website's performance with professional speed optimization services. I analyze and optimize your site to reduce loading times and enhance user experience, leading to higher rankings in search engines.",
		"Showcase your products with stunning photography that highlights every detail. Whether it's for an online store, catalog, or social media, I provide high-resolution images that make your products look their best.",
		"Get noticed with a digital business card designed to leave a lasting impression. I create interactive, shareable cards that reflect your personal or business brand, perfect for networking in the digital age.",
		"Enhance customer interactions with a custom AI-powered chatbot. I develop chatbots that understand user intent and provide instant, accurate responses, improving customer satisfaction and efficiency.",
		"Make your podcast sound professional with expert editing services. I remove background noise, adjust audio levels, and add effects, ensuring your episodes are clear and engaging.",
		"Stand out in the job market with a professionally written resume and cover letter. I craft documents that highlight your skills and achievements, increasing your chances of landing your dream job.",
		"Capture attention with interactive web banners designed to engage users. I create banners with dynamic elements that encourage clicks and improve your website's conversion rates.",
		"Visualize complex information with custom infographic designs. I transform data into visually appealing graphics that are easy to understand and share, perfect for presentations and marketing.",
		"Create a comprehensive online course that educates and engages your audience. I handle everything from content creation to video production, ensuring a high-quality learning experience.",
		"Streamline your tasks with reliable data entry and virtual assistance services. I provide accurate data management and administrative support, allowing you to focus on more critical aspects of your business.",
		"Optimize your LinkedIn profile to attract recruiters and potential clients. I enhance your profile's visibility by crafting a compelling summary, highlighting key achievements, and adding relevant skills.",
		"Explain complex ideas with animated explainer videos that simplify concepts. I create videos that are engaging, informative, and tailored to your target audience, perfect for marketing and educational purposes.",
		"Enhance your photos with professional retouching and editing. Whether it's portraits, product photos, or event shots, I ensure your images look their best with adjustments to color, lighting, and details.",
		"Grow your online presence with expert social media management. I create, schedule, and monitor posts across various platforms, ensuring consistent engagement and brand visibility.",
		"Extend your WordPress site's functionality with a custom plugin. I develop plugins tailored to your specific needs, providing features that enhance user experience and streamline operations.",
		"Unlock the full potential of your data with advanced Excel analysis. I provide insights through data visualization, statistical analysis, and automated reports, helping you make informed decisions.",
	}
	categories := []string{
		"Graphic Design",
		"Web Development",
		"Writing & Translation",
		"Digital Marketing",
		"Video & Animation",
		"Programming & Tech",
		"Music & Audio",
		"Business",
		"Lifestyle",
		"Art & Design",
		"Photography",
		"Marketing Strategy",
		"Online Education",
		"Administrative Support",
		"Data Analysis",
		"Consulting",
		"App Development",
		"Customer Service",
		"Legal Services",
		"Financial Planning",
	}
	subCategoryList := []string{
		"Logo Design",
		"WordPress",
		"SEO Writing",
		"Mobile App Design",
		"Social Media Marketing",
		"E-commerce",
		"Voiceover",
		"Video Editing",
		"Illustration",
		"Business Card Design",
		"T-Shirt Design",
		"Email Marketing",
		"3D Rendering",
		"Copywriting",
		"Speed Optimization",
		"Product Photography",
		"Digital Business Cards",
		"AI Chatbot Development",
		"Podcast Editing",
		"Resume Writing",
	}
	tagList := []string{
		"Design", "Logo", "Branding",
		"WordPress", "Website", "Development",
		"SEO", "Content", "Writing",
		"UI", "UX", "Mobile App",
		"Social Media", "Marketing", "Strategy",
		"E-commerce", "Online Store", "Shopify",
		"Voiceover", "Commercial", "Narration",
		"Video Editing", "YouTube", "Production",
		"Illustration", "Art", "Drawing",
		"Business Card", "Graphic Design", "Print",
		"3D", "Rendering", "Animation",
		"T-Shirt", "Merchandise", "Custom",
		"Email", "Campaign", "Optimization",
		"Photography", "Product", "Photoshoot",
		"Chatbot", "AI", "Automation",
		"Podcast", "Audio", "Editing",
		"Resume", "Job Application", "Career",
		"Infographic", "Visual", "Data",
		"Excel", "Data Analysis", "Automation",
		"Consulting", "Strategy", "Business",
		"App Development", "iOS", "Android",
		"Customer Service", "Support", "CRM",
		"Legal", "Contracts", "Advice",
		"Finance", "Planning", "Investment",
	}
	getRandomItem := func(items []string) string {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		return items[rand.Intn(len(items))]
	}

	for i := 0; i < 1000; i++ {
		title := getRandomItem(titles)
		description := getRandomItem(descriptions)
		category := getRandomItem(categories)
		subCategories := pq.StringArray{getRandomItem(subCategoryList), getRandomItem(subCategoryList), getRandomItem(subCategoryList), getRandomItem(subCategoryList), getRandomItem(subCategoryList)}
		tags := pq.StringArray{getRandomItem(tagList), getRandomItem(tagList), getRandomItem(tagList), getRandomItem(tagList), getRandomItem(tagList)}
		randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		five, four, three, two, one := randomizer.Intn(10)*randomizer.Intn(10), randomizer.Intn(10)*randomizer.Intn(10), randomizer.Intn(10)*randomizer.Intn(10), randomizer.Intn(10)*randomizer.Intn(10), randomizer.Intn(10)*randomizer.Intn(10)
		result := db.Exec(
			`INSERT INTO gigs (id, seller_id, title, description, category, sub_categories, tags, active,
			expected_delivery_days, ratings_count, rating_sum, rating_categories, price, cover_image, created_at,
			title_tokens, description_tokens, category_tokens, sub_categories_tokens, tags_tokens)
			VALUES (
			uuid_generate_v4(),
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?::jsonb,
			?,
			?,
			?,
			strip(to_tsvector('english', ?)),
			strip(to_tsvector('english', ?)),
			strip(to_tsvector('english', ?)),
			strip(to_tsvector('english', ?)),
			strip(to_tsvector('english', ?)))`,
			"tlGIviGcfMSPIdkOkssCg6nOPZPOqLozMOgRSRVBDNlyBwQwNhAElMPBXMWg65mG",
			title,
			description,
			category,
			subCategories,
			tags,
			true,                            // active
			randomizer.Intn(10)*3,           // expected delivery
			uint64(five+four+three+two+one), // rating count
			uint64(5*five+4*four+3*three+2*two+one), // rating sum
			// rating categories
			fmt.Sprintf(`{"five": %d,"four": %d,"three": %d,"two": %d,"one": %d}`,
				five, four, three, two, one),
			randomizer.Intn(10)*randomizer.Intn(10), // price
			faker.URL(), // cover image
			time.Now(),  // created at
			title,
			description,
			category,
			strings.Join(subCategories, ","),
			strings.Join(tags, ","),
			// subCategories,
			// tags,
		)

		if result.Error != nil {
			fmt.Printf("seeding data error:\n%+v", result.Error)
			continue
		}
	}
}
