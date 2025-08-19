package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Contact represents a contact form submission
type Contact struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Naam       string    `json:"naam" gorm:"not null"`
	Bedrijf    string    `json:"bedrijf"`
	Email      string    `json:"email" gorm:"not null"`
	Telefoon   string    `json:"telefoon"`
	Onderwerp  string    `json:"onderwerp" gorm:"not null"`
	Urgentie   string    `json:"urgentie"`
	Bericht    string    `json:"bericht" gorm:"not null"`
	Privacy    bool      `json:"privacy" gorm:"not null"`
	Nieuwsbrief bool     `json:"nieuwsbrief"`
	CreatedAt  time.Time `json:"created_at"`
}

// PageData represents data passed to templates
type PageData struct {
	Title       string
	Description string
	Page        string
	Content     template.HTML
}

var db *gorm.DB
var geminiClient *genai.GenerativeModel

func main() {
	// Initialize database
	initDatabase()

	// Initialize Gemini client
	initGeminiClient()

	// Create Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// Routes
	r.GET("/", homeHandler)
	r.GET("/diensten", dienstenHandler)
	r.GET("/over-ons", overOnsHandler)
	r.GET("/contact", contactGetHandler)
	r.POST("/contact", contactPostHandler)
	r.GET("/privacybeleid", privacybeleidHandler)

	// New route for Gemini chat
	r.POST("/chat", chatHandler)

	// Start server
	r.Run("0.0.0.0:8080")
}

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("ict_eerbeek.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Contact{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}

func initGeminiClient() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	geminiClient = client.GenerativeModel("gemini-pro")
}

func chatHandler(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	resp, err := geminiClient.GenerateContent(ctx, genai.Text(request.Message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		response := resp.Candidates[0].Content.Parts[0].(genai.Text)
		c.JSON(http.StatusOK, gin.H{"reply": string(response)})
	} else {
		c.JSON(http.StatusOK, gin.H{"reply": "No response from AI."})
	}
}

func homeHandler(c *gin.Context) {
	data := PageData{
		Title:       "Home",
		Description: "ICT Eerbeek - Uw betrouwbare partner voor alle ICT-oplossingen in Eerbeek en omgeving. Netwerk & security, website ontwerp, IoT & AI oplossingen, en computerhulp.",
		Page:        "home",
	}

	// Read home template content
	homeContent := `<!-- Hero Section -->
<section class="hero">
    <div class="hero-container">
        <h1>Welkom bij ICT Eerbeek</h1>
        <p>Uw betrouwbare partner voor alle ICT-oplossingen in Eerbeek en omgeving. Van netwerk & security tot website ontwerp en AI-oplossingen.</p>
        <a href="/contact" class="cta-button">Neem Contact Op</a>
    </div>
</section>

<!-- Services Section -->
<section class="services">
    <div class="services-container">
        <h2 class="section-title">Onze Diensten</h2>
        <div class="services-grid">
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-shield-alt"></i>
                </div>
                <h3>Netwerk & Security</h3>
                <p>Professionele netwerkoplossingen en beveiligingssystemen om uw bedrijf te beschermen tegen cyberdreigingen en optimale prestaties te garanderen.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-palette"></i>
                </div>
                <h3>Website & Logo Ontwerp</h3>
                <p>Creatieve en professionele website- en logo-ontwerpen die uw merk versterken en uw online aanwezigheid verbeteren.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-microchip"></i>
                </div>
                <h3>IoT & AI Oplossingen</h3>
                <p>Innovatieve Internet of Things en Artificial Intelligence oplossingen om uw bedrijfsprocessen te automatiseren en optimaliseren.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-tools"></i>
                </div>
                <h3>All-round Computerhulp</h3>
                <p>Uitgebreide computerondersteuning voor particulieren en bedrijven, van hardware reparaties tot software installaties en training.</p>
            </div>
        </div>
    </div>
</section>

<!-- About Preview Section -->
<section class="about-preview" style="padding: 6rem 0; background: linear-gradient(135deg, #f8f9fa, #e9ecef);">
    <div class="services-container">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 4rem; align-items: center;">
            <div>
                <h2 style="font-size: 2.5rem; font-weight: 700; margin-bottom: 1.5rem; color: var(--text-dark);">Waarom ICT Eerbeek?</h2>
                <p style="font-size: 1.1rem; line-height: 1.7; color: var(--text-light); margin-bottom: 2rem;">
                    Met jarenlange ervaring in de ICT-sector bieden wij betrouwbare en innovatieve oplossingen voor al uw technologische uitdagingen. 
                    Ons team van experts staat klaar om uw bedrijf naar het volgende niveau te tillen.
                </p>
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                    <div style="text-align: center; padding: 1rem;">
                        <div style="font-size: 2rem; font-weight: 700; color: var(--primary-green);">10+</div>
                        <div style="color: var(--text-light);">Jaar Ervaring</div>
                    </div>
                    <div style="text-align: center; padding: 1rem;">
                        <div style="font-size: 2rem; font-weight: 700; color: var(--primary-blue);">100+</div>
                        <div style="color: var(--text-light);">Tevreden Klanten</div>
                    </div>
                    <div style="text-align: center; padding: 1rem;">
                        <div style="font-size: 2rem; font-weight: 700; color: var(--primary-purple);">24/7</div>
                        <div style="color: var(--text-light);">Ondersteuning</div>
                    </div>
                    <div style="text-align: center; padding: 1rem;">
                        <div style="font-size: 2rem; font-weight: 700; color: var(--primary-brown);">100%</div>
                        <div style="color: var(--text-light);">Tevredenheid</div>
                    </div>
                </div>
            </div>
            <div style="text-align: center;">
                <img src="/static/images/logo.jpg" alt="ICT Eerbeek Team" style="max-width: 100%; height: auto; border-radius: 15px; box-shadow: 0 10px 30px rgba(0,0,0,0.1);">
            </div>
        </div>
    </div>
</section>

<!-- Contact CTA Section -->
<section style="background: linear-gradient(135deg, var(--primary-green), var(--primary-blue)); color: white; padding: 4rem 0; text-align: center;">
    <div class="services-container">
        <h2 style="font-size: 2.5rem; font-weight: 700; margin-bottom: 1rem;">Klaar om te beginnen?</h2>
        <p style="font-size: 1.2rem; margin-bottom: 2rem; opacity: 0.9;">
            Neem vandaag nog contact met ons op voor een vrijblijvende consultatie en ontdek hoe wij uw ICT-uitdagingen kunnen oplossen.
        </p>
        <div style="display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap;">
            <a href="/contact" class="cta-button">Contact Opnemen</a>
            <a href="/diensten" style="display: inline-block; background: transparent; color: white; padding: 1rem 2rem; text-decoration: none; border: 2px solid white; border-radius: 50px; font-weight: 600; transition: all 0.3s ease;">Bekijk Diensten</a>
        </div>
    </div>
</section>`

	data.Content = template.HTML(homeContent)
	c.HTML(http.StatusOK, "base.html", data)
}

func dienstenHandler(c *gin.Context) {
	data := PageData{
		Title:       "Onze Diensten",
		Description: "Ontdek ons uitgebreide aanbod van ICT-oplossingen: netwerk & security, website & logo ontwerp, IoT & AI oplossingen, en all-round computerhulp.",
		Page:        "diensten",
	}

	dienstenContent := `<!-- Page Header -->
<section class="page-header">
    <div class="hero-container">
        <h1>Onze Diensten</h1>
        <p>Ontdek ons uitgebreide aanbod van ICT-oplossingen voor particulieren en bedrijven</p>
    </div>
</section>

<!-- Page Content -->
<div class="page-content">
    <!-- Netwerk & Security -->
    <div class="content-section" id="netwerk-security">
        <h2><i class="fas fa-shield-alt" style="color: var(--primary-blue); margin-right: 1rem;"></i>Netwerk & Security</h2>
        <p>In de digitale wereld van vandaag is een betrouwbaar netwerk en sterke beveiliging essentieel voor elk bedrijf. Wij bieden uitgebreide netwerkoplossingen en beveiligingsdiensten om uw bedrijf te beschermen.</p>
        
        <div class="services-grid" style="margin-top: 2rem;">
            <div class="service-card">
                <h3>Netwerkinstallatie</h3>
                <p>Professionele installatie van bedrijfsnetwerken, inclusief bekabeling, switches, routers en access points voor optimale connectiviteit.</p>
            </div>
            <div class="service-card">
                <h3>Firewall Configuratie</h3>
                <p>Implementatie en configuratie van geavanceerde firewalls om uw netwerk te beschermen tegen externe bedreigingen.</p>
            </div>
            <div class="service-card">
                <h3>VPN Oplossingen</h3>
                <p>Veilige externe toegang tot uw bedrijfsnetwerk via Virtual Private Network oplossingen voor thuiswerkers.</p>
            </div>
            <div class="service-card">
                <h3>Security Monitoring</h3>
                <p>24/7 monitoring van uw netwerk om verdachte activiteiten te detecteren en direct actie te ondernemen.</p>
            </div>
        </div>
    </div>

    <!-- Website & Logo Ontwerp -->
    <div class="content-section" id="website-logo">
        <h2><i class="fas fa-palette" style="color: var(--primary-green); margin-right: 1rem;"></i>Website & Logo Ontwerp</h2>
        <p>Uw online aanwezigheid is cruciaal voor het succes van uw bedrijf. Wij creëren professionele websites en memorabele logo's die uw merk versterken en klanten aantrekken.</p>
        
        <div class="services-grid" style="margin-top: 2rem;">
            <div class="service-card">
                <h3>Responsive Webdesign</h3>
                <p>Moderne websites die perfect werken op alle apparaten, van desktop tot smartphone, met focus op gebruikerservaring.</p>
            </div>
            <div class="service-card">
                <h3>E-commerce Oplossingen</h3>
                <p>Volledige webshops met betalingssystemen, voorraadbeheersystemen en klantenportalen voor online verkoop.</p>
            </div>
            <div class="service-card">
                <h3>Logo & Branding</h3>
                <p>Creatieve logo-ontwerpen en complete huisstijlen die uw bedrijf onderscheiden van de concurrentie.</p>
            </div>
            <div class="service-card">
                <h3>SEO Optimalisatie</h3>
                <p>Zoekmachine optimalisatie om uw website beter vindbaar te maken in Google en andere zoekmachines.</p>
            </div>
        </div>
    </div>

    <!-- IoT & AI Oplossingen -->
    <div class="content-section" id="iot-ai">
        <h2><i class="fas fa-microchip" style="color: var(--primary-purple); margin-right: 1rem;"></i>IoT & AI Oplossingen</h2>
        <p>Stap in de toekomst met onze innovatieve Internet of Things en Artificial Intelligence oplossingen. Automatiseer processen, verzamel waardevolle data en optimaliseer uw bedrijfsvoering.</p>
        
        <div class="services-grid" style="margin-top: 2rem;">
            <div class="service-card">
                <h3>Smart Building Systemen</h3>
                <p>Intelligente gebouwbeheersystemen voor verlichting, klimaatbeheersing en beveiliging met IoT-sensoren.</p>
            </div>
            <div class="service-card">
                <h3>Industriële Automatisering</h3>
                <p>IoT-oplossingen voor productieprocessen, kwaliteitscontrole en voorspellend onderhoud in de industrie.</p>
            </div>
            <div class="service-card">
                <h3>AI Chatbots</h3>
                <p>Intelligente chatbots voor klantenservice die 24/7 beschikbaar zijn en veel voorkomende vragen automatisch beantwoorden.</p>
            </div>
            <div class="service-card">
                <h3>Data Analytics</h3>
                <p>AI-gedreven data-analyse om patronen te herkennen, trends te voorspellen en betere bedrijfsbeslissingen te nemen.</p>
            </div>
        </div>
    </div>

    <!-- All-round Computerhulp -->
    <div class="content-section" id="computerhulp">
        <h2><i class="fas fa-tools" style="color: var(--primary-brown); margin-right: 1rem;"></i>All-round Computerhulp</h2>
        <p>Van hardware reparaties tot software installaties, wij bieden uitgebreide computerondersteuning voor particulieren en bedrijven. Geen probleem is te klein of te groot.</p>
        
        <div class="services-grid" style="margin-top: 2rem;">
            <div class="service-card">
                <h3>Hardware Reparatie</h3>
                <p>Reparatie van computers, laptops, printers en andere hardware met snelle diagnose en eerlijke prijzen.</p>
            </div>
            <div class="service-card">
                <h3>Software Installatie</h3>
                <p>Installatie en configuratie van besturingssystemen, applicaties en drivers voor optimale prestaties.</p>
            </div>
            <div class="service-card">
                <h3>Data Recovery</h3>
                <p>Herstel van verloren data van harde schijven, USB-sticks en andere opslagmedia met geavanceerde technieken.</p>
            </div>
            <div class="service-card">
                <h3>IT Training</h3>
                <p>Persoonlijke training en workshops om uw digitale vaardigheden te verbeteren en efficiënter te werken.</p>
            </div>
        </div>
    </div>

    <!-- Call to Action -->
    <div class="highlight-box">
        <h3>Heeft u vragen over onze diensten?</h3>
        <p>Neem contact met ons op voor een vrijblijvende consultatie. Wij denken graag met u mee over de beste ICT-oplossing voor uw situatie.</p>
        <a href="/contact" class="cta-button" style="margin-top: 1rem; display: inline-block;">Contact Opnemen</a>
    </div>
</div>`

	data.Content = template.HTML(dienstenContent)
	c.HTML(http.StatusOK, "base.html", data)
}

func overOnsHandler(c *gin.Context) {
	data := PageData{
		Title:       "Over ICT Eerbeek",
		Description: "Leer meer over ICT Eerbeek, ons team, onze missie en onze passie voor technologie. Uw betrouwbare ICT-partner in Eerbeek.",
		Page:        "over-ons",
	}

	overOnsContent := `<!-- Page Header -->
<section class="page-header">
    <div class="hero-container">
        <h1>Over ICT Eerbeek</h1>
        <p>Leer meer over ons team, onze missie en onze passie voor technologie</p>
    </div>
</section>

<!-- Page Content -->
<div class="page-content">
    <!-- Company Story -->
    <div class="content-section">
        <h2>Ons Verhaal</h2>
        <p>ICT Eerbeek is ontstaan uit de passie voor technologie en de wens om bedrijven en particulieren te helpen bij hun digitale uitdagingen. Sinds onze oprichting hebben wij ons ontwikkeld tot een betrouwbare partner voor alle ICT-gerelateerde vraagstukken in Eerbeek en omgeving.</p>
        
        <p>Wat begon als een kleine onderneming is uitgegroeid tot een professioneel ICT-bedrijf dat zich onderscheidt door persoonlijke service, technische expertise en innovatieve oplossingen. Wij geloven dat technologie toegankelijk moet zijn voor iedereen, ongeacht de grootte van uw bedrijf of uw technische achtergrond.</p>
    </div>

    <!-- Mission & Vision -->
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 3rem; margin: 3rem 0;">
        <div class="service-card">
            <div class="service-icon">
                <i class="fas fa-bullseye"></i>
            </div>
            <h3>Onze Missie</h3>
            <p>Wij maken technologie toegankelijk en begrijpelijk voor iedereen. Door persoonlijke service en maatwerkoplossingen helpen wij onze klanten hun digitale doelen te bereiken en hun bedrijfsprocessen te optimaliseren.</p>
        </div>
        <div class="service-card">
            <div class="service-icon">
                <i class="fas fa-eye"></i>
            </div>
            <h3>Onze Visie</h3>
            <p>Wij streven ernaar de meest vertrouwde en innovatieve ICT-partner te zijn, die duurzame oplossingen levert die bijdragen aan het succes van onze klanten in een steeds digitalere wereld.</p>
        </div>
    </div>

    <!-- Team Section -->
    <div class="content-section">
        <h2>Ons Team</h2>
        <p>Ons team bestaat uit gepassioneerde en gecertificeerde ICT-professionals met jarenlange ervaring in diverse vakgebieden. Wij werken nauw samen om de beste oplossingen te leveren en staan altijd klaar om u te ondersteunen.</p>
        
        <div class="team-grid">
            <div class="team-member">
                <img src="https://via.placeholder.com/100" alt="Teamlid 1">
                <h4>Jan de Vries</h4>
                <p>Oprichter & Lead Netwerk Engineer</p>
            </div>
            <div class="team-member">
                <img src="https://via.placeholder.com/100" alt="Teamlid 2">
                <h4>Sophie Jansen</h4>
                <p>Webdesigner & UI/UX Specialist</p>
            </div>
            <div class="team-member">
                <img src="https://via.placeholder.com/100" alt="Teamlid 3">
                <h4>Mark van Dijk</h4>
                <p>IoT & AI Ontwikkelaar</p>
            </div>
            <div class="team-member">
                <img src="https://via.placeholder.com/100" alt="Teamlid 4">
                <h4>Linda Bakker</h4>
                <p>All-round IT Support Specialist</p>
            </div>
        </div>
    </div>

    <!-- Values Section -->
    <div class="content-section">
        <h2>Onze Waarden</h2>
        <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 2rem; margin-top: 2rem;">
            <div class="service-card">
                <h3>Klantgerichtheid</h3>
                <p>De klant staat centraal in alles wat we doen. Wij luisteren naar uw behoeften en leveren oplossingen die echt waarde toevoegen.</p>
            </div>
            <div class="service-card">
                <h3>Innovatie</h3>
                <p>Wij blijven op de hoogte van de nieuwste technologische ontwikkelingen en passen deze toe om u de meest geavanceerde oplossingen te bieden.</p>
            </div>
            <div class="service-card">
                <h3>Betrouwbaarheid</h3>
                <p>U kunt op ons rekenen. Wij leveren wat we beloven en zorgen voor stabiele en veilige ICT-omgevingen.</p>
            </div>
        </div>
    </div>

    <!-- Call to Action -->
    <div class="highlight-box">
        <h3>Benieuwd wat wij voor u kunnen betekenen?</h3>
        <p>Neem contact op voor een vrijblijvend gesprek. Wij helpen u graag verder!</p>
        <a href="/contact" class="cta-button" style="margin-top: 1rem; display: inline-block;">Contact Opnemen</a>
    </div>
</div>`

	data.Content = template.HTML(overOnsContent)
	c.HTML(http.StatusOK, "base.html", data)
}

func contactGetHandler(c *gin.Context) {
	data := PageData{
		Title:       "Contact",
		Description: "Neem contact op met ICT Eerbeek voor al uw vragen over netwerk & security, website ontwerp, IoT & AI oplossingen, en computerhulp.",
		Page:        "contact",
	}
	c.HTML(http.StatusOK, "contact.html", data)
}

func contactPostHandler(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact.CreatedAt = time.Now()

	if result := db.Create(&contact); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bericht succesvol verzonden!"})
}

func privacybeleidHandler(c *gin.Context) {
	data := PageData{
		Title:       "Privacybeleid",
		Description: "Lees het privacybeleid van ICT Eerbeek. Wij respecteren uw privacy en zorgen voor een veilige verwerking van uw persoonsgegevens.",
		Page:        "privacybeleid",
	}
	c.HTML(http.StatusOK, "privacybeleid.html", data)
}


