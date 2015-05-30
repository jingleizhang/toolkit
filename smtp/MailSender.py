#!/usr/bin/python
# coding: utf-8
'''

usage：
cat body.txt | ./MailSender.py [-h] -s SUBJECT -o SENDER -p PASSWD -d DESTINATION 
                     [-c CC] [-b BCC] [-H] [-f [FILES [FILES ...]]]

optional arguments:
  -h, --help            show this help message and exit
  -s SUBJECT, --subject SUBJECT
                        email Subject
  -o SENDER, --sender SENDER
                        email Sender
  -p PASSWD, --passwd PASSWD
                        email Sender's password
  -d DESTINATION, --destination DESTINATION
                        email to ,
  -c CC, --cc CC        CC to ,
  -b BCC, --bcc BCC     BCC to ,
  -H                    TEXT_HTML
  -f [FILES [FILES ...]], --files [FILES [FILES ...]]
                        Switch of attachments
'''
import time
import logging
import os
import sys
import subprocess
import argparse
import chardet
import smtplib
# For guessing MIME type based on file name extension
import mimetypes
from email import encoders
from email.message import Message
from email.mime.audio import MIMEAudio
from email.mime.base import MIMEBase
from email.mime.image import MIMEImage
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.Utils import COMMASPACE, formatdate

def parse_argvs(argvs):
	parser = argparse.ArgumentParser(description="MailSender")
	parser.add_argument("-s", '--subject', help="email Subject", required=True, action="store")
	parser.add_argument("-o", '--sender',  help="email Sender", required=True, action="store")
	parser.add_argument("-p", '--passwd',  help="email Sender's password", action="store")
	parser.add_argument("-d", '--destination', help="email to,", required=True, action="store")
	parser.add_argument("-c", '--cc', help="CC to ,", action="store")
	parser.add_argument("-b", '--bcc', help="BCC to ,", action="store")
	parser.add_argument("-H", help="TEXT_HTML", action='store_true', default=False)
	parser.add_argument('-f', '--files', nargs='*', help="Switch of attachments", action="store")   
	args = parser.parse_args(argvs)
	return args

def to_unicode(str):
	'''
	Trans data into unicode.
	'''
	charResult = chardet.detect(str)
	for key in charResult:
		if key == "confidence":
			continue
		#print charResult[key]
		try:
			str = str.decode(charResult[key])
		except:
			continue
		break
	return str

if __name__ == '__main__':
	args = parse_argvs(sys.argv[1:])

	LOGGER = logging.getLogger('Mail.Result')  
	LOGGER.setLevel(logging.INFO)  

	fh = logging.FileHandler('MailResult.log', 'w')  
	fh.setLevel(logging.INFO)  

	ch = logging.StreamHandler()  
	ch.setLevel(logging.INFO)  

	formatter = logging.Formatter('[%(asctime)s] [%(levelname)s] [%(message)s]')  
	fh.setFormatter(formatter)  
	ch.setFormatter(formatter)  

	LOGGER.addHandler(fh)  
	LOGGER.addHandler(ch) 

	SERVER = 'xxx'

	USER = args.sender
	PASSWD = args.passwd


	mail = smtplib.SMTP(SERVER)
	try:
		mail.ehlo()
		#mail.login(USER, PASSWD)
		LOGGER.info("finish login")		

		msg = MIMEMultipart()
		msg['Subject'] = to_unicode(args.subject).encode('utf-8')					
        msg['From'] = args.sender + ","

		toList = ""						
		toTmp = args.destination.split(",")
		toTmp = sorted(set(toTmp), key=toTmp.index)
		for towho in toTmp:
			if towho != "":
                		toList += towho + ";"
		msg['To'] = toList


		if args.cc is not None:
			ccList = ""				
			ccTmp = args.cc.split(";")
			ccTmp = sorted(set(ccTmp), key=ccTmp.index)
			for ccwho in ccTmp:
				if ccwho != "":
                    			ccList += ccwho + ","
			msg['Cc'] = ccList

		if args.bcc is not None:
			bccList = ""				
			bccTmp = args.bcc.split(";")
			bccTmp = sorted(set(bccTmp), key=bccTmp.index)
			for bccwho in bccTmp:
				if bccwho != "":
                    			bccList += bccwho + ","
			msg['Bcc'] = bccList

		body = sys.stdin.read()
		body = to_unicode(body).encode('utf-8')
		if args.H:
			# html
			msg.attach(MIMEText(body, 'html', 'utf-8'))
		else:
			# plain
			msg.attach(MIMEText(body, 'plain', 'utf-8'))

		if args.files is not None:
			for furl in args.files:
				if os.path.exists(furl):
					# Guess the content type based on the file's extension.  Encoding
					# will be ignored, although we should check for simple things like
					# gzip'd or compressed files.
					ctype, encoding = mimetypes.guess_type(furl)
					if ctype is None or encoding is not None:
						# No guess could be made, or the file is encoded (compressed), so
						# use a generic bag-of-bits type.
						ctype = 'application/octet-stream'
					maintype, subtype = ctype.split('/', 1)
					if maintype == 'text':
						fp = open(furl)
						# Note: we should handle calculating the charset
						msg_attach = MIMEText(fp.read(), _subtype=subtype)
						fp.close()
					elif maintype == 'image':
						fp = open(furl, 'rb')
						msg_attach = MIMEImage(fp.read(), _subtype=subtype)
						fp.close()
					elif maintype == 'audio':
						fp = open(furl, 'rb')
						msg_attach = MIMEAudio(fp.read(), _subtype=subtype)
						fp.close()
					else:
						fp = open(furl, 'rb')
						msg_attach = MIMEBase(maintype, subtype)
						msg_attach.set_payload(fp.read())
						fp.close()
						# Encode the payload using Base64
						encoders.encode_base64(msg_attach)
					# Set the filename parameter
                    			msg_attach.add_header('Content-ID', '<%s>' % os.path.basename(furl))
					#msg_attach.add_header('Content-Disposition', 'attachment', filename=os.path.basename(furl))
                    			msg_attach.add_header('Content-Disposition', 'inline', filename=os.path.basename(furl))
					msg.attach(msg_attach)
				else:
					LOGGER.info("file %s not exist !!!" %(furl)) 


		msg['Date'] = formatdate(localtime=True)

		mail_to = msg['To']
		if msg['Cc'] is not None:
			mail_to += msg['Cc']
		if msg['Bcc'] is not None:
			mail_to += msg['Bcc']

		mail.sendmail(msg['From'], mail_to.split(","), msg.as_string())
		LOGGER.info("finish send mail")

	except smtplib.SMTPAuthenticationError:
		LOGGER.error("SMTPAuthenticationError user: %s passwd: %s" %(USER, PASSWD))
		sys.exit(-1)
	except smtplib.SMTPSenderRefused, ex:
		LOGGER.error("SMTPSenderRefused: %s" %(ex))
		sys.exit(-2)
	except smtplib.SMTPRecipientsRefused, ex:
		LOGGER.error("SMTPRecipientsRefused: %s" %(ex))
		sys.exit(-3)
	except smtplib.SMTPDataError, ex:
		LOGGER.error("SMTPDataError: %s" %(ex))
		sys.exit(-4)

	LOGGER.info("send mail success. from: %s to: %s subject: %s" %(msg['From'], msg['To'], msg['Subject']))
	mail.quit()
	sys.exit(0)


