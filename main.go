package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func main() {
	// Initialize database
	initDatabase()

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
            <p>Wij streven ernaar de meest vertrouwde ICT-partner in de regio te zijn, bekend om onze innovatieve oplossingen, uitstekende service en langdurige klantrelaties gebaseerd op vertrouwen en expertise.</p>
        </div>
    </div>

    <!-- Values -->
    <div class="content-section">
        <h2>Onze Waarden</h2>
        <div class="services-grid">
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-handshake"></i>
                </div>
                <h3>Betrouwbaarheid</h3>
                <p>Wij staan voor onze afspraken en leveren altijd wat wij beloven. Uw vertrouwen is ons belangrijkste bezit.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-lightbulb"></i>
                </div>
                <h3>Innovatie</h3>
                <p>Wij blijven op de hoogte van de nieuwste technologieën en trends om u de beste oplossingen te kunnen bieden.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-users"></i>
                </div>
                <h3>Persoonlijke Service</h3>
                <p>Elke klant is uniek. Wij luisteren naar uw specifieke behoeften en bieden maatwerkoplossingen.</p>
            </div>
            <div class="service-card">
                <div class="service-icon">
                    <i class="fas fa-graduation-cap"></i>
                </div>
                <h3>Kennis Delen</h3>
                <p>Wij geloven in het delen van kennis en helpen onze klanten om zelfstandiger te worden met technologie.</p>
            </div>
        </div>
    </div>

    <!-- Team Section -->
    <div class="content-section">
        <h2>Ons Team</h2>
        <p>Ons team bestaat uit ervaren ICT-professionals die elk hun eigen specialisatie hebben. Samen vormen wij een sterk team dat klaar staat om al uw ICT-uitdagingen aan te gaan.</p>
        
        <div class="team-grid">
            <div class="team-member">
                <img src="/static/images/logo.jpg" alt="Team Member">
                <h4>Jan van der Berg</h4>
                <p>Oprichter & Senior ICT Consultant</p>
                <p style="margin-top: 1rem; font-size: 0.9rem;">Specialist in netwerkbeveiliging en systeembeheer met meer dan 15 jaar ervaring in de ICT-sector.</p>
            </div>
            <div class="team-member">
                <img src="/static/images/logo.jpg" alt="Team Member">
                <h4>Sarah Jansen</h4>
                <p>Web Developer & Designer</p>
                <p style="margin-top: 1rem; font-size: 0.9rem;">Creatieve webdesigner en developer gespecialiseerd in moderne websites en gebruikerservaring.</p>
            </div>
            <div class="team-member">
                <img src="/static/images/logo.jpg" alt="Team Member">
                <h4>Mike de Vries</h4>
                <p>IoT & AI Specialist</p>
                <p style="margin-top: 1rem; font-size: 0.9rem;">Expert in Internet of Things en Artificial Intelligence oplossingen voor bedrijfsautomatisering.</p>
            </div>
        </div>
    </div>

    <!-- Call to Action -->
    <div class="content-section" style="text-align: center;">
        <h2>Klaar om Samen te Werken?</h2>
        <p>Wij kijken ernaar uit om kennis te maken en te horen hoe wij u kunnen helpen met uw ICT-uitdagingen.</p>
        <a href="/contact" class="cta-button" style="margin-top: 2rem; display: inline-block;">Neem Contact Op</a>
    </div>
</div>`

	data.Content = template.HTML(overOnsContent)
	c.HTML(http.StatusOK, "base.html", data)
}

func contactGetHandler(c *gin.Context) {
	data := PageData{
		Title:       "Contact",
		Description: "Neem contact op met ICT Eerbeek voor al uw ICT-vragen en uitdagingen. Wij staan klaar om u te helpen.",
		Page:        "contact",
	}

	contactContent := `<!-- Page Header -->
<section class="page-header">
    <div class="hero-container">
        <h1>Contact</h1>
        <p>Neem contact met ons op voor al uw ICT-vragen en uitdagingen</p>
    </div>
</section>

<!-- Page Content -->
<div class="page-content">
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 4rem; align-items: start;">
        <!-- Contact Information -->
        <div>
            <h2>Contactgegevens</h2>
            <p style="margin-bottom: 2rem;">Wij staan klaar om u te helpen. Neem gerust contact met ons op via onderstaande gegevens of vul het contactformulier in.</p>
            
            <div class="contact-info" style="margin-bottom: 3rem;">
                <div style="display: flex; align-items: center; margin-bottom: 1rem; padding: 1rem; background: var(--light-gray); border-radius: 10px;">
                    <i class="fas fa-envelope" style="color: var(--primary-blue); font-size: 1.5rem; margin-right: 1rem; width: 30px;"></i>
                    <div>
                        <strong>E-mail</strong><br>
                        <a href="mailto:info@ict-eerbeek.nl" style="color: var(--primary-blue); text-decoration: none;">info@ict-eerbeek.nl</a>
                    </div>
                </div>
                
                <div style="display: flex; align-items: center; margin-bottom: 1rem; padding: 1rem; background: var(--light-gray); border-radius: 10px;">
                    <i class="fas fa-phone" style="color: var(--primary-green); font-size: 1.5rem; margin-right: 1rem; width: 30px;"></i>
                    <div>
                        <strong>Telefoon</strong><br>
                        <a href="tel:+31612345678" style="color: var(--primary-green); text-decoration: none;">+31 (0)6 12345678</a>
                    </div>
                </div>
                
                <div style="display: flex; align-items: center; margin-bottom: 1rem; padding: 1rem; background: var(--light-gray); border-radius: 10px;">
                    <i class="fas fa-map-marker-alt" style="color: var(--primary-purple); font-size: 1.5rem; margin-right: 1rem; width: 30px;"></i>
                    <div>
                        <strong>Locatie</strong><br>
                        Eerbeek, Nederland
                    </div>
                </div>
                
                <div style="display: flex; align-items: center; margin-bottom: 1rem; padding: 1rem; background: var(--light-gray); border-radius: 10px;">
                    <i class="fas fa-clock" style="color: var(--primary-brown); font-size: 1.5rem; margin-right: 1rem; width: 30px;"></i>
                    <div>
                        <strong>Openingstijden</strong><br>
                        Ma-Vr: 09:00 - 17:00<br>
                        Za: 10:00 - 14:00<br>
                        Zo: Gesloten
                    </div>
                </div>
            </div>
            
            <h3>Spoedgevallen</h3>
            <p>Voor urgente ICT-problemen zijn wij 24/7 bereikbaar via ons spoednummer:</p>
            <div style="background: linear-gradient(135deg, var(--primary-green), var(--light-green)); color: white; padding: 1.5rem; border-radius: 10px; text-align: center; margin-top: 1rem;">
                <i class="fas fa-phone" style="font-size: 1.5rem; margin-bottom: 0.5rem;"></i><br>
                <strong style="font-size: 1.2rem;">+31 (0)6 87654321</strong><br>
                <small>24/7 Spoednummer</small>
            </div>
        </div>
        
        <!-- Contact Form -->
        <div>
            <h2>Stuur ons een bericht</h2>
            <form id="contact-form" class="contact-form">
                <div class="form-group">
                    <label for="naam">Naam *</label>
                    <input type="text" id="naam" name="naam" required>
                </div>
                
                <div class="form-group">
                    <label for="bedrijf">Bedrijf</label>
                    <input type="text" id="bedrijf" name="bedrijf">
                </div>
                
                <div class="form-group">
                    <label for="email">E-mailadres *</label>
                    <input type="email" id="email" name="email" required>
                </div>
                
                <div class="form-group">
                    <label for="telefoon">Telefoonnummer</label>
                    <input type="tel" id="telefoon" name="telefoon">
                </div>
                
                <div class="form-group">
                    <label for="onderwerp">Onderwerp *</label>
                    <select id="onderwerp" name="onderwerp" required>
                        <option value="">Selecteer een onderwerp</option>
                        <option value="netwerk-security">Netwerk & Security</option>
                        <option value="website-logo">Website & Logo Ontwerp</option>
                        <option value="iot-ai">IoT & AI Oplossingen</option>
                        <option value="computerhulp">All-round Computerhulp</option>
                        <option value="offerte">Offerte Aanvraag</option>
                        <option value="ondersteuning">Technische Ondersteuning</option>
                        <option value="anders">Anders</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="urgentie">Urgentie</label>
                    <select id="urgentie" name="urgentie">
                        <option value="laag">Laag - Binnen een week</option>
                        <option value="normaal" selected>Normaal - Binnen 2-3 dagen</option>
                        <option value="hoog">Hoog - Binnen 24 uur</option>
                        <option value="urgent">Urgent - Zo spoedig mogelijk</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="bericht">Bericht *</label>
                    <textarea id="bericht" name="bericht" placeholder="Beschrijf uw vraag of probleem zo gedetailleerd mogelijk..." required></textarea>
                </div>
                
                <div class="form-group">
                    <label style="display: flex; align-items: center; cursor: pointer;">
                        <input type="checkbox" name="privacy" required style="margin-right: 0.5rem;">
                        Ik ga akkoord met het <a href="/privacybeleid" style="color: var(--primary-blue);">privacybeleid</a> *
                    </label>
                </div>
                
                <div class="form-group">
                    <label style="display: flex; align-items: center; cursor: pointer;">
                        <input type="checkbox" name="nieuwsbrief" style="margin-right: 0.5rem;">
                        Ik wil graag op de hoogte blijven van nieuws en aanbiedingen
                    </label>
                </div>
                
                <button type="submit" class="submit-button">
                    <i class="fas fa-paper-plane" style="margin-right: 0.5rem;"></i>
                    Bericht Versturen
                </button>
            </form>
        </div>
    </div>
    
    <!-- FAQ Section -->
    <div class="content-section" style="margin-top: 4rem;">
        <h2>Veelgestelde Vragen</h2>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-top: 2rem;">
            <div class="service-card">
                <h3>Hoe snel krijg ik antwoord?</h3>
                <p>Wij streven ernaar om binnen 4 uur te reageren op werkdagen. Voor urgente zaken kunt u ons spoednummer bellen.</p>
            </div>
            <div class="service-card">
                <h3>Zijn er kosten verbonden aan een consult?</h3>
                <p>Het eerste consult en de offerte zijn altijd gratis. Pas na goedkeuring van de offerte brengen wij kosten in rekening.</p>
            </div>
            <div class="service-card">
                <h3>Werken jullie ook in het weekend?</h3>
                <p>Voor spoedgevallen zijn wij 24/7 bereikbaar. Reguliere werkzaamheden plannen wij in overleg, ook in het weekend indien gewenst.</p>
            </div>
            <div class="service-card">
                <h3>Bieden jullie onderhoud contracten?</h3>
                <p>Ja, wij bieden verschillende onderhoudscontracten aan voor zowel particulieren als bedrijven. Neem contact op voor meer informatie.</p>
            </div>
        </div>
    </div>
</div>`

	data.Content = template.HTML(contactContent)
	c.HTML(http.StatusOK, "base.html", data)
}

func contactPostHandler(c *gin.Context) {
	var contact Contact

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid form data"})
		return
	}

	// Validate required fields
	if contact.Naam == "" || contact.Email == "" || contact.Onderwerp == "" || contact.Bericht == "" || !contact.Privacy {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Missing required fields"})
		return
	}

	// Set default urgency if not provided
	if contact.Urgentie == "" {
		contact.Urgentie = "normaal"
	}

	// Set creation time
	contact.CreatedAt = time.Now()

	// Save to database
	if err := db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to save contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Contact form submitted successfully"})
}

func privacybeleidHandler(c *gin.Context) {
	data := PageData{
		Title:       "Privacybeleid",
		Description: "Privacybeleid van ICT Eerbeek - Hoe wij omgaan met uw persoonlijke gegevens en welke rechten u heeft.",
		Page:        "privacybeleid",
	}

	privacyContent := `<!-- Page Header -->
<section class="page-header">
    <div class="hero-container">
        <h1>Privacybeleid</h1>
        <p>Hoe wij omgaan met uw persoonlijke gegevens</p>
    </div>
</section>

<!-- Page Content -->
<div class="page-content">
    <div class="content-section">
        <p><strong>Laatst bijgewerkt:</strong> 1 januari 2024</p>
        
        <h2>1. Inleiding</h2>
        <p>ICT Eerbeek hecht grote waarde aan de bescherming van uw persoonlijke gegevens. In dit privacybeleid leggen wij uit welke persoonlijke gegevens wij verzamelen, hoe wij deze gebruiken en welke rechten u heeft met betrekking tot uw gegevens.</p>
        
        <h2>2. Contactgegevens</h2>
        <div class="service-card">
            <h3>Verantwoordelijke voor de gegevensverwerking:</h3>
            <p>
                <strong>ICT Eerbeek</strong><br>
                E-mail: info@ict-eerbeek.nl<br>
                Telefoon: +31 (0)6 12345678<br>
                Adres: Eerbeek, Nederland
            </p>
        </div>
        
        <h2>3. Welke gegevens verzamelen wij?</h2>
        <p>Wij kunnen de volgende categorieën persoonlijke gegevens van u verzamelen:</p>
        
        <h3>3.1 Contactgegevens</h3>
        <ul style="margin-left: 2rem; margin-bottom: 1rem;">
            <li>Naam en achternaam</li>
            <li>E-mailadres</li>
            <li>Telefoonnummer</li>
            <li>Bedrijfsnaam (indien van toepassing)</li>
            <li>Adresgegevens</li>
        </ul>
        
        <h2>4. Contact en klachten</h2>
        <div class="highlight-box">
            <h3>Contact opnemen</h3>
            <p>
                <strong>E-mail:</strong> privacy@ict-eerbeek.nl<br>
                <strong>Telefoon:</strong> +31 (0)6 12345678<br>
                <strong>Post:</strong> ICT Eerbeek, Eerbeek, Nederland
            </p>
        </div>
    </div>
</div>`

	data.Content = template.HTML(privacyContent)
	c.HTML(http.StatusOK, "base.html", data)
}

