import Head from 'next/head'
import Navbar from './components/navbar'
import HomePage from './components/home'
import Footer from './components/footer'
import AddDisease from './components/addDisease'
import TestDNA from './components/testDNA'
import SearchBar from './components/searchbar'
import styles from '../styles/Home.module.css'
import AboutUs from './components/aboutus'
import ScrollTop from './components/scroll'

export default function Home() {
  return (
    <div>
      <Head>
        <title>DNA Pattern Matching</title>
        <link rel="preload" href="\font\Montserrat\static\Montserrat-Thin.ttf" as="font" crossOrigin="" />
        <link rel="preload" href="\font\Montserrat\static\Montserrat-Light.ttf" as="font" crossOrigin="" />v
        <link rel="preload" href="\font\Montserrat\static\Montserrat-Regular.ttf" as="font" crossOrigin="" />
        <link rel="preload" href="\font\Montserrat\static\Montserrat-Bold.ttf" as="font" crossOrigin="" />
        <meta name="description" content= "Tugas Besar 3 STIMA" />

        <meta name="title" content="DNA Pattern Matching" />
        <meta itemProp="description" content="Tugas Besar 3 STIMA" />
        <meta itemProp="image" content="https://drive.google.com/uc?export=view&id=1ZhhiWAX97xFsF1sL0bxyH4E23aJ24H-F" />
        <link rel="icon" href="https://drive.google.com/uc?export=view&id=1ZhhiWAX97xFsF1sL0bxyH4E23aJ24H-F" />
        <link rel="apple-touch-icon" href="https://drive.google.com/uc?export=view&id=1ZhhiWAX97xFsF1sL0bxyH4E23aJ24H-F" />

        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </Head>
      <Navbar />
      <div id="HomeSection">
        <HomePage />
      </div>
      <div id="AddDiseaseSection">
        <AddDisease />
      </div>
      <div id="TestDNASection">
        <TestDNA />
      </div>
      <div id="SearchBarSection">
        <SearchBar />
      </div>
      <div id="AboutUsSection">
        <AboutUs />
      </div>
      <div id="ContactSection">
        <Footer />
      </div>
      <ScrollTop />
    </div>
  )
}
