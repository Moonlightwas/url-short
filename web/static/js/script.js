document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById("url-form")
    const urlInput = document.getElementById('url-input');
    const shortButton = document.getElementById('short-button');
    const resultDiv = document.getElementById('result-div');
    const resultUrl = document.getElementById('result-url');
    const copyButton = document.getElementById('copy-button');
    const canvas = document.getElementById('qr-code-canvas');
    const downloadPNG = document.getElementById('download-qr-png')
    const downloadSVG = document.getElementById('download-qr-svg')


    form.addEventListener('submit', async (e) => {
        e.preventDefault(); 
        await handleShortButtonClick();
    });

    async function generateQR(canvasElement, url, options = {
        width: 300,
        margin: 2,
        color: {
            dark: '#f8fafc',
            light: '#0f172a'
        },
    }) {
        return new Promise((resolve, reject) => {
            QRCode.toCanvas(canvasElement, url, options, (err) => {
                if (err) {
                    reject(err);
                } else {
                    resolve();
                }
            });
        });
    };

    shortButton.addEventListener('click', async function () {
        handleShortButtonClick();
    });

    async function handleShortButtonClick() {
        try {
            const url = urlInput.value.trim()
            if (!url){
                alert("URL can't be empty");
                return 
            }

            if (!url.startsWith("https://") && !url.startsWith("http://")){
                alert("Field must contain a url");
                return
            }

            const response = await fetch('/', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({url: url}),
            });

            if (!response.ok){
                const errorText = await response.text();
                throw new Error(errorText);
            }

            const data = await response.json();
            resultUrl.textContent = data.processedURL;

            await generateQR(canvas, resultUrl.textContent)

            resultDiv.style.display = 'block';
        } catch (err) {
            console.error('Error:', err);
            alert(err.message);
        }
    };

    copyButton.addEventListener('click', function() {
        const textToCopy = document.getElementById('result-url').textContent;
        const tempInput = document.createElement('textarea');
        try {
            tempInput.value = textToCopy; 
            document.body.appendChild(tempInput);
            tempInput.select();
            document.execCommand('copy');
            document.body.removeChild(tempInput);

            alert('URL copied to clipboard!');
        } catch(err) {
            console.error('Error:', err);
            alert('Failed to copy');
        }
        window.getSelection().removeAllRanges();
    });

    downloadPNG.addEventListener('click', function(){
        downloadQR('png')
    })

    downloadSVG.addEventListener('click', function(){
        downloadQR('svg')
    })

    async function downloadQR(format) {
        try {
            if (format == 'png') {
                const link = document.createElement('a');
                link.download = 'url-short.png';
                link.href = canvas.toDataURL('image/png');
                link.click();
            } else if (format == 'svg') {
                const svgData = await QRCode.toString(resultUrl.textContent, {
                    type: 'svg',
                    margin: 2,
                    color: {
                        dark: '#f8fafc',
                        light: '#0f172a'
                    }
                });
            
                const blob = new Blob([svgData], { type: 'image/svg+xml' });
                const svgUrl = URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.download = 'url-short.svg';
                link.href = svgUrl;
                link.click();
            }
        } catch(err) {
            console.error('Error:', err);
            alert(err.message);
        }
    };
});
