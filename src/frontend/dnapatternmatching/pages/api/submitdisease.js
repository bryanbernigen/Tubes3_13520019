export default function handler(req, res) {
    if (req.method === 'POST') {
      const namapenyakit = req.body.namapenyakit
      const rantaidna = req.body.rantaidna
      const penyakit = {
        namapenyakit,
        rantaidna
      }
      res.status(200).json(penyakit)
    }
  }